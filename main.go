package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/kelseyhightower/envconfig"
	"github.com/rcrowley/go-metrics"
	"github.com/therealbill/libredis/client"
)

type User struct {
	Name string
	ID   int
}

type Game struct {
	Name string
	ID   int
}

type Match struct {
	Player1 User
	Player2 User
	GameID  Game
	Winner  string
	Loser   string
}

type LaunchConfig struct {
	Name                  string
	RedisConnectionString string
	RedisAuthToken        string
	JSONOut               bool
	ReaderCount           int
	UserCount             int
	MatchesPerUser        int
	TotalMatches          int
	GameServerCount       int
	PoolSize              int
	Pipeline              bool
	GameName              string
}

var (
	ds                    *client.Redis
	initial_command_time  float64
	post_command_time     float64
	initial_command_count float64
	post_command_count    float64
	initial_command_times map[string]float64
	commands_tested       []string
	TTT                   Game
	Users                 []User
	resc                  chan int
	config                LaunchConfig
	report                bytes.Buffer
	m                     runtime.MemStats
)

func init() {
	err := envconfig.Process("commissar", &config)
	if config.UserCount == 0 {
		config.UserCount = 10
	}
	if config.ReaderCount == 0 {
		config.ReaderCount = 1
	}
	if config.PoolSize == 0 {
		config.PoolSize = 1
	}
	if config.MatchesPerUser == 0 {
		config.MatchesPerUser = 100
	}
	if config.GameServerCount == 0 {
		config.GameServerCount = 1
	}
	if config.GameName == "" {
		config.GameName = "Tic-Tac-Toe"
	}
	if config.RedisConnectionString == "" {
		config.RedisConnectionString = "127.0.0.1:6379"
	}

	config.TotalMatches = config.MatchesPerUser * config.UserCount

	TTT = Game{Name: config.GameName, ID: 1}
	Users = make([]User, config.UserCount)
	for x := 1; x < config.UserCount+1; x++ {
		player := User{Name: fmt.Sprintf("Player-%d", x), ID: x}
		Users[x-1] = player
	}
	initial_command_times = make(map[string]float64)
	s := metrics.NewUniformSample(config.TotalMatches)
	h := metrics.NewHistogram(s)
	metrics.Register("latency:full", h)
	commands_tested = []string{"zincrby", "hincrby", "zrevrange"}
	ds, err = client.DialWithConfig(&client.DialConfig{Address: config.RedisConnectionString, Password: config.RedisAuthToken, MaxIdle: config.PoolSize})
	if err != nil {
		log.Fatal("Unable to connect to redis:", err)
	}
	ds.FlushAll()
	if ds == nil {
		log.Fatal("No redis connection, no testing")
	}
	cmdstats, err := ds.Info()
	if err != nil {
		log.Fatalf("Err: %s", err.Error())
	}
	for _, c := range commands_tested {
		for k, v := range cmdstats.Commandstats.Stats {
			if k == c {
				initial_command_time += v["usec"]
				initial_command_times[k] = v["usec"]
				initial_command_count += v["calls"]
			}
		}
	}
	resc = make(chan int)
}

func serveGames(matches int) {
	dsG, err := client.DialWithConfig(&client.DialConfig{Address: config.RedisConnectionString, Password: config.RedisAuthToken})
	if err != nil {
		log.Fatal("Redis Connection failure:", err)
	}
	for m := 1; m <= matches; m++ {
		pid1 := rand.Intn(config.UserCount)
		pid2 := rand.Intn(config.UserCount)
		if pid2 == pid1 {
			pid2 = rand.Intn(config.UserCount)
		}
		match := Match{Player1: Users[pid1], Player2: Users[pid2], GameID: TTT}
		playGame(&match)
		win_userstatkey := fmt.Sprintf("users:%s:%s", match.Winner, TTT.Name)
		loss_userstatkey := fmt.Sprintf("users:%s:%s", match.Loser, TTT.Name)
		lbkey := "leaderboard:" + TTT.Name
		updateStats(dsG, win_userstatkey, loss_userstatkey, lbkey, &match)
	}
	resc <- 1
}

func playGame(m *Match) *Match {
	rand.Seed(time.Now().UnixNano())
	winner := rand.Intn(1)
	switch winner {
	case 0:
		m.Winner = m.Player1.Name
		m.Loser = m.Player2.Name
	case 1:
		m.Winner = m.Player2.Name
		m.Loser = m.Player1.Name
	}
	//time.Sleep(time.Duration(int64(rand.Intn(20))) * time.Millisecond)
	return m
}

func updateStats(ds *client.Redis, win_userstatkey, loss_userstatkey, lbkey string, match *Match) {
	if config.Pipeline {
		p, err := ds.Pipelining()
		if err != nil {
			fmt.Print("Error obtaining pipeline")
		}
		defer p.Close()
		p.Command("HINCRBY", win_userstatkey, "netwins", 1)
		p.Command("HINCRBY", win_userstatkey, "wins", 1)
		p.Command("HINCRBY", loss_userstatkey, "loss", 1)
		p.Command("HINCRBY", loss_userstatkey, "netwins", -1)
		p.Command("ZINCRBY", lbkey, 1.0, match.Winner)
		p.Command("ZINCRBY", lbkey, -1.0, match.Loser)
		start := time.Now() // this is where we actually send the data so we should track from here
		p.ReceiveAll()
		elapsed := time.Duration(time.Since(start) / 6)
		h := metrics.Get("latency:full").(metrics.Histogram)
		h.Update(int64(elapsed))

	} else {
		start := time.Now()
		ds.HIncrBy(win_userstatkey, "wins", 1)
		elapsed := time.Since(start)
		h := metrics.Get("latency:full").(metrics.Histogram)
		h.Update(int64(elapsed))
		start = time.Now()
		ds.HIncrBy(win_userstatkey, "netwins", 1)
		elapsed = time.Since(start)
		h = metrics.Get("latency:full").(metrics.Histogram)
		h.Update(int64(elapsed))
		start = time.Now()
		ds.HIncrBy(loss_userstatkey, "loss", 1)
		elapsed = time.Since(start)
		h = metrics.Get("latency:full").(metrics.Histogram)
		h.Update(int64(elapsed))
		start = time.Now()
		ds.HIncrBy(loss_userstatkey, "netwins", -1)
		elapsed = time.Since(start)
		h = metrics.Get("latency:full").(metrics.Histogram)
		h.Update(int64(elapsed))
		start = time.Now()
		ds.ZIncrBy(lbkey, 1.0, match.Winner)
		elapsed = time.Since(start)
		h = metrics.Get("latency:full").(metrics.Histogram)
		h.Update(int64(elapsed))
		start = time.Now()
		ds.ZIncrBy(lbkey, -1.0, match.Loser)
		elapsed = time.Since(start)
		h = metrics.Get("latency:full").(metrics.Histogram)
		h.Update(int64(elapsed))
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	h := metrics.Get("latency:full").(metrics.Histogram)
	h.Update(int64(elapsed))
}

func pullStats() {
	rds, err := client.DialWithConfig(&client.DialConfig{Address: config.RedisConnectionString, Password: config.RedisAuthToken})
	if err != nil {
		log.Fatalf("puller unable to connect: %s", err.Error())
	}
	lbkey := "leaderboard:" + TTT.Name
	for {
		start := time.Now()
		rds.ZRevRange(lbkey, 0, 10, true)
		elapsed := time.Since(start)
		h := metrics.Get("latency:full").(metrics.Histogram)
		h.Update(int64(elapsed))
		sleep := time.Duration(int64(rand.Intn(20))) * time.Millisecond
		time.Sleep(sleep)
	}
}

func main() {
	for r := 0; r < config.ReaderCount; r++ {
		go pullStats()
	}
	for x := 0; x < config.GameServerCount; x++ {
		go serveGames(config.TotalMatches / config.GameServerCount)
	}
	for i := 0; i < config.GameServerCount; i++ {
		select {
		case res := <-resc:
			_ = res
		}
	}
	runtime.ReadMemStats(&m)
	fmt.Printf("%d,%d,%d,%d\n", m.HeapSys, m.HeapAlloc, m.HeapIdle, m.HeapReleased)
	h := metrics.Get("latency:full").(metrics.Histogram)
	snap := h.Snapshot()
	cmdstats, err := ds.Info()
	if err != nil {
		log.Fatalf("Err: %s", err.Error())
	}
	for _, c := range commands_tested {
		for k, v := range cmdstats.Commandstats.Stats {
			if k == c {
				report.WriteString(fmt.Sprintf("Server-side test time for %s = %v\n", k, time.Duration(v["usec"]-initial_command_times[k])))
				post_command_time += v["usec"]
				post_command_count += v["calls"]
			}
		}
	}
	//report.WriteString(fmt.Sprintf("post_command_count: %.0f\n", post_command_count))
	//report.WriteString(fmt.Sprintf("initial_command_count: %.0f\n", initial_command_count))

	config.RedisAuthToken = ""
	b, _ := json.Marshal(config)
	report.WriteString(string(b))
	commands_processed := int64(post_command_count - initial_command_count)
	report.WriteString(fmt.Sprintf("\n%d commands_processed in %v\n", commands_processed, time.Duration(snap.Sum())))
	num_seconds := float64(snap.Sum()) / float64(time.Second)
	cps := float64(commands_processed) / num_seconds
	gps := float64(config.TotalMatches) / num_seconds
	report.WriteString(fmt.Sprintf("# %s games simulated in %v\n", humanize.Comma(int64(config.TotalMatches)), time.Duration(snap.Sum())))
	report.WriteString(fmt.Sprintf("Games/second:   \t%s\n", humanize.Comma(int64(gps))))
	report.WriteString(fmt.Sprintf("Commands/second:\t%s\n", humanize.Comma(int64(cps))))
	buckets := []float64{0.99, 0.95, 0.9, 0.75, 0.5}
	dist := snap.Percentiles(buckets)
	min := time.Duration(snap.Min())
	max := time.Duration(snap.Max())
	mean := time.Duration(snap.Mean())
	stddev := time.Duration(snap.StdDev())
	report.WriteString(fmt.Sprintf("Max: %s\nMin: %s\nMean: %s\nJitter: %s\n", max, min, mean, stddev))
	report.WriteString("\nPercentile breakout:")
	report.WriteString("\n====================\n")
	for i, b := range buckets {
		d := time.Duration(dist[i])
		report.WriteString(fmt.Sprintf("%.2f%%: %v\n", b*100, d))
	}
	server_time := int64(post_command_time - initial_command_time)
	log.Printf("snap.Sum(): %v\n", snap.Sum())
	total_time := int64(snap.Sum())
	net_time := total_time - server_time
	server_pct := float64(server_time) / float64(total_time)
	network_pct := float64(net_time) / float64(total_time)
	report.WriteString(fmt.Sprintf("Total client time: %v\n", time.Duration(total_time)))
	report.WriteString(fmt.Sprintf("Server time: %v (%2.2f%%)\n", time.Duration(server_time), server_pct*100))
	report.WriteString(fmt.Sprintf("Network/Client time: %v (%2.2f%%)\n", time.Duration(net_time), network_pct*100))
	print(report.String())
}

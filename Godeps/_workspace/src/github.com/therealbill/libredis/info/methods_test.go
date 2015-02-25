package info

import (
	"fmt"
	"strings"
	"testing"
	"time"

	rclient "github.com/therealbill/libredis/client"
)

var (
	network  = "tcp"
	address  = "127.0.0.1:6379"
	db       = 1
	password = ""
	timeout  = 5 * time.Second
	maxidle  = 1
	r        *rclient.Redis
	format   = "tcp://auth:%s@%s/%d?timeout=%s&maxidle=%d"
)

func init() {
	client, err := rclient.DialWithConfig(&rclient.DialConfig{network, address, db, password, timeout, maxidle})
	if err != nil {
		panic(err)
	}
	r = client
}

func TestBuildMapFromInfoString(t *testing.T) {
	raw, _ := r.InfoString("server")
	trimmed := strings.TrimSpace(raw)
	rmap := BuildMapFromInfoString(trimmed)
	if len(rmap["redis_version"]) == 0 {
		t.Error("Version wasn't parsed")
		t.Fail()
	}
}

func TestBuildInfoKeyspace(t *testing.T) {
	keyinfo, _ := r.InfoString("keyspace")
	space := BuildInfoKeyspace(keyinfo)
	if len(space.Databases) == 0 {
		t.Fail()
	}
}

func TestCommandStats(t *testing.T) {
	res, _ := r.InfoString("commandstats")
	stats := CommandStats(res)
	if len(stats.Stats) == 0 {
		t.Fail()
	}
	// The act of calling CommandStats will produce at least one call it info
	// So we ensure we have at least one call
	if stats.Stats["info"]["calls"] == 0 {
		t.Fail()
	}
}

func TestKeyspaceStats(t *testing.T) {
	r.Set("deletme", "1")
	res, _ := r.InfoString("keyspace")
	stats := KeyspaceStats(res)
	if len(stats.Databases) == 0 {
		fmt.Printf("%+v\n", stats)
		t.Error("didn't get expected Keyspace Stats structure.")
		t.Fail()
	}
	r.Del("deleteme")
}

func TestBuildAllInfoMap(t *testing.T) {
	res, _ := r.InfoString("all")
	alldata := BuildAllInfoMap(res)
	if len(alldata["CPU"]["used_cpu_sys"]) == 0 {
		fmt.Printf("alldata.cpu.used_cpu_sys: %+v\n", alldata["CPU"]["used_cpu_sys"])
		t.Error("didn't parse cpu.used_cpu_sys")
		t.Fail()
	}
}

func TestGetAllInfo(t *testing.T) {
	res, _ := r.InfoString("all")
	all := GetAllInfo(res)
	// Server Checks
	if all.Server.Arch_bits == 0 {
		t.Error("Didn't parse Server.Arch_bits")
		t.Fail()
	}
}

func TestInfo(t *testing.T) {
	r.Set("deleteme", "1")
	all, err := r.Info()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	// Server Checks
	if all.Server.Arch_bits == 0 {
		t.Error("Didn't parse Server.Arch_bits")
		t.Fail()
	}
	if len(all.Server.Version) == 0 {
		t.Error("Didn't parse Server.Version")
		t.Fail()
	}
	// Replication Checks
	if len(all.Replication.Role) == 0 {
		t.Error("Failed to parse Replication.Role")
		t.Fail()
	}

	// Persistence
	if all.Persistence.AOFEnabled {
		t.Error("Tests assume default config, so this AOFEnabled shoudl be false")
		t.Fail()
	}
	// Stats
	if all.Stats.TotalCommandsProcessed == 0 {
		t.Error("Failed to parse Stats.TotalCommandsProcessed")
		t.Fail()
	}
	// Memory
	if len(all.Memory.UsedMemoryHuman) == 0 {
		t.Error("Failed to parse Memory.UsedMemoryHuman")
		t.Fail()
	}
	if all.Memory.UsedMemory == 0 {
		t.Error("Failed to parse Memory.UsedMemory")
		t.Fail()
	}
	// Keyspace
	if len(all.Keyspace.Databases) == 0 {
		t.Error("Failed to parse at least one DB from Keyspace")
		t.Fail()
	}
	if all.Commandstats.Stats["info"]["calls"] == 0 {
		t.Error("Failed to parse stats on info command from Commandstats")
		fmt.Printf("Commandstats:\t%+v\n", all.Commandstats)
		t.Fail()
	}
	r.Del("deleteme", "1")

}

func TestUpperFirst(t *testing.T) {
	instring := "this"
	outsting := upperFirst(instring)
	if instring == outsting {
		t.Error("Failed to convert this to This")
		t.Fail()
	}
	if outsting != "This" {
		t.Error("Failed to convert this to This")
		t.Fail()
	}

	empty := upperFirst("")
	if empty != "" {
		t.Error("upperFirst on empty strign result sin non-empty string")
		t.Fail()
	}
}

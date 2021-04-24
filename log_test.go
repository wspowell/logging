package logging

import (
	"sync"
	"testing"
)

func Test_Log_GlobalTags_race(t *testing.T) {
	cfg := newTestConfig(LevelError)
	log := NewLog(cfg)

	wg := &sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			log.Error("test")
		}()
	}

	wg.Wait()
}

func Test_Log_Tags(t *testing.T) {
	cfg := newTestConfig(LevelError)
	cfg.Tags()["global"] = "global"

	log := NewLog(cfg)

	log.Tag("test1", "value1")
	log.Tag("test2", "value2")

	if value, ok := log.cfg.Tags()["global"]; !ok || value != "global" {
		t.Error("missing global tag")
	}

	if value, ok := log.tags["test1"]; !ok || value != "value1" {
		t.Error("missing test1 tag")
	}

	if value, ok := log.tags["test2"]; !ok || value != "value2" {
		t.Error("missing test2 tag")
	}
}

func Test_Log_Tags_nil(t *testing.T) {
	cfg := newTestConfig(LevelError)
	log := NewLog(cfg)

	log.Tag("test1", "value1")
	log.Tag("test2", "value2")

	if value, ok := log.tags["test1"]; !ok || value != "value1" {
		t.Fail()
	}

	if value, ok := log.tags["test2"]; !ok || value != "value2" {
		t.Fail()
	}
}

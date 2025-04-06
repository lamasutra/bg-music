package model

import (
	"regexp"
	"strings"
	"time"

	"github.com/lamasutra/bg-music/wt-client/client"
	"github.com/lamasutra/bg-music/wt-client/ui"
)

type Player struct {
	Name                    string
	Vehicle                 string
	Dead                    bool
	Damaged                 bool
	SeverlyDamaged          bool
	LastKillTime            int64
	LastDamageTime          int64
	LastSeverDamageTime     int64
	LastBurnTime            int64
	LastKilledTime          int64
	LastDamagedTime         int64
	LastSeverelyDamagedTime int64
	LastBurnedTime          int64
	Targets                 map[string]*Player
	CurrentTarget           *Player
}

type DamageParser struct {
	regexp          *regexp.Regexp
	player          map[string]*Player
	lastAnyKillTime int64
}

func NewDamageParser() *DamageParser {
	d := DamageParser{}
	d.init()

	return &d
}

func (d *DamageParser) init() {
	d.regexp = regexp.MustCompile(`(\[ai\])?\s?([^\(]+)\s+(\([^\)]+\))?\s?(critically\sdamaged|severely\sdamaged|shot\sdown|destroyed|has\scrashed\.|set\safire)\s*(\[ai\])?\s?([^\(]+)?\s?(\([^\)]+\))?\.?`)
	d.player = make(map[string]*Player, 20)
}

func (d *DamageParser) Parse(hudMsg *client.HudMsg) {
	hudMsg.Each(d.parseDamage)
}

func (d *DamageParser) GetLastKillTime() int64 {
	return d.lastAnyKillTime
}

func (d *DamageParser) parseDamage(dmg client.Damage, index int) bool {
	ui.Debug("parse damage", dmg.Msg)
	matches := d.regexp.FindStringSubmatch(dmg.Msg)
	if len(matches) == 0 {
		ui.Debug("no match")
		return false
	}

	i := 0
	for i < len(matches) {
		ui.Debug("match", i, matches[i])
		i++
	}
	// ui.Debug(matches)
	var aiSource, aiTarget bool
	aiSource = matches[1] == "[ai]"
	aiTarget = matches[5] == "[ai]" || (matches[5] == "" && matches[7] == "")
	if aiSource {
		ui.Debug("ai source")
		d.parseAiSource(&matches)
	} else if aiTarget {
		ui.Debug("ai target")
		d.parseAiTarget(&matches)
	} else {
		ui.Debug("players")
		d.parsePlayers(&matches)
	}

	return false
}

func (d *DamageParser) parseAiSource(matches *[]string) {
	// var sourceName, targetName string
	// var sourceVehicle, targetVehicle string
	// var action string

}

func (d *DamageParser) parseAiTarget(matches *[]string) (*Player, *Player) {
	var sourceName, targetName string
	var sourceVehicle, targetVehicle string
	var action string

	sourceName = strings.TrimRight((*matches)[2], " ")
	sourceVehicle = (*matches)[3]
	action = (*matches)[4]
	targetName = strings.TrimRight((*matches)[5], " ")
	targetVehicle = (*matches)[6]

	sourcePlayer := d.FindOrCreatePlayer(sourceName)
	targetPlayer := d.FindOrCreatePlayer(targetName)
	sourcePlayer.Vehicle = strings.TrimRight(strings.TrimLeft(sourceVehicle, "("), ")")
	targetPlayer.Vehicle = strings.TrimRight(strings.TrimLeft(targetVehicle, "("), ")")

	// ui.Debug("source:", sourceName, "sourceVehicle:", sourceVehicle, "action:", action, "target:", targetName, "targetVehicle:", targetVehicle)
	handled := d.handleAction(action, sourcePlayer, targetPlayer)
	if handled {
		ui.Debug("action", action, "handled")
	} else {
		ui.Debug("action", action, "not handled")
	}

	return sourcePlayer, targetPlayer
}

func (d *DamageParser) parsePlayers(matches *[]string) (*Player, *Player) {
	var sourceName, targetName string
	var sourceVehicle, targetVehicle string
	var action string

	sourceName = strings.TrimRight((*matches)[2], " ")
	sourceVehicle = (*matches)[3]
	action = (*matches)[4]
	targetName = strings.TrimRight((*matches)[6], " ")
	targetVehicle = (*matches)[7]

	sourcePlayer := d.FindOrCreatePlayer(sourceName)
	targetPlayer := d.FindOrCreatePlayer(targetName)
	sourcePlayer.Vehicle = strings.TrimRight(strings.TrimLeft(sourceVehicle, "("), ")")
	targetPlayer.Vehicle = strings.TrimRight(strings.TrimLeft(targetVehicle, "("), ")")

	// ui.Debug("source:", sourceName, "sourceVehicle:", sourceVehicle, "action:", action, "target:", targetName, "targetVehicle:", targetVehicle)
	handled := d.handleAction(action, sourcePlayer, targetPlayer)
	if handled {
		ui.Debug("action", action, "handled")
	} else {
		ui.Debug("action", action, "not handled")
	}

	return sourcePlayer, targetPlayer
}

func (d *DamageParser) FindOrCreatePlayer(name string) *Player {
	pl, ok := d.player[name]
	if ok {
		ui.Debug("player found", name)
		return pl
	}
	ui.Debug("creating player", name)
	pl = &Player{
		Name:    name,
		Targets: make(map[string]*Player, 64),
	}

	d.player[name] = pl

	return pl
}

func (d *DamageParser) handleAction(action string, source *Player, target *Player) bool {
	now := time.Now().Unix()
	switch action {
	case "has crashed.":
		source.Damaged = true
		source.Dead = true
		source.SeverlyDamaged = true
		source.LastKilledTime = now
		source.LastDamagedTime = now
		source.LastSeverelyDamagedTime = now
		ui.Debug("has crashed", source, target)
		return true
	case "shot down", "destroyed":
		source.LastKillTime = now
		target.LastKilledTime = now
		target.Dead = true
		target.Damaged = true
		source.addTarget(target)
		ui.Debug("shot down or destroyed", source, target)
		d.lastAnyKillTime = now
		return true
	case "set afire":
		source.LastBurnTime = now
		target.LastBurnedTime = now
		target.Damaged = true
		source.addTarget(target)
		ui.Debug("set afire", source, target)
		return true
	case "critically damaged":
		source.LastDamageTime = now
		target.LastDamagedTime = now
		target.Damaged = true
		source.addTarget(target)
		ui.Debug("critically damaged", source, target)
		return true
	case "severely damaged":
		source.LastSeverDamageTime = now
		target.LastSeverelyDamagedTime = now
		target.Damaged = true
		source.addTarget(target)
		ui.Debug("severely damaged", source, target)
		return true
	}

	return false
}

func (p *Player) addTarget(t *Player) {
	_, ok := p.Targets[t.Name]
	if !ok {
		p.Targets[t.Name] = t
	}
	p.CurrentTarget = t
}

func (p *Player) GetTarget(name string) *Player {
	t, ok := p.Targets[name]
	if ok {
		return t
	}

	return nil
}

func (p *Player) Reset() {
	p.Vehicle = ""
	p.CurrentTarget = nil
	p.Damaged = false
	p.Dead = false
	p.LastBurnTime = 0
	p.LastBurnedTime = 0
	p.LastDamageTime = 0
	p.LastDamagedTime = 0
	p.LastKillTime = 0
	p.LastKilledTime = 0
	p.LastSeverDamageTime = 0
	p.LastSeverelyDamagedTime = 0
	p.Targets = make(map[string]*Player, 32)
}

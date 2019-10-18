package main

import "sort"

type League []Player

func (l League) Find(name string) *Player {
	for i, player := range l {
		if player.Name == name {
			return &l[i]
		}
	}
	return nil
}

func (l League) Sort(sortBy func(i, j int)bool) {
	sort.Slice(l, sortBy)
}

func SortByPlayerWins(league League) func(i, j int) bool {
	return func(i, j int) bool {
		return league[i].Wins > league[j].Wins
	}
}

func SortByPlayerName(league League) func(i, j int) bool {
	return func(i, j int) bool {
		return league[i].Name < league[j].Name
	}
}


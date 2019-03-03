package utils

func (s *Shuffle) Shuffle(userIDs []string) []string {
	for i := len(userIDs); i > 0; i-- {
		randIndex := s.Rand.Intn(i)
		userIDs[i-1], userIDs[randIndex] = userIDs[randIndex], userIDs[i-1]
	}
	return userIDs
}

func (s *Shuffle) Split(userIDs []string) [][]string {
	var groups [][]string
	var stragglers []string
	stdGroupCnt := len(userIDs) / standardGroupSize
	stdGroupUserCnt := standardGroupSize * stdGroupCnt
	if len(userIDs) > stdGroupUserCnt {
		for _, strag := range userIDs[stdGroupUserCnt:] {
			stragglers = append(stragglers, strag)
		}
	}

	for i := 0; i < stdGroupCnt; i++ {
		var group []string
		for i := 0; i < standardGroupSize; i++ {
			group = append(group, userIDs[0])
			userIDs = append(userIDs[:0], userIDs[0+1:]...)
		}
		groups = append(groups, group)
	}

	if len(groups) == 0 {
		return nil
	}

	if stragglers != nil {
		for i, strag := range stragglers {
			groups[i] = append(groups[i], strag)
		}
	}

	return groups
}

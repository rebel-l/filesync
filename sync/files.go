package sync

type Files []File

func (f Files) SpaceNeeded() int64 {
	var space int64

	for _, v := range f {
		space += v.Source.Info.Size() // TODO: space of destination file (if exists) must be subtracted for correct number
	}

	return space
}

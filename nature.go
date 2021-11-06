package bippa

type Nature string

func (nature Nature) IsValid() bool {
	_, ok := NATUREDEX[nature]
	return ok
}

type Natures []Nature

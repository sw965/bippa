package bippa

import (
	"fmt"
)


//関数の上の行のコメントは図鑑ナンバー

//3
func NewVenusaur() Pokemon {
	pokemon, err := NewPokemon("フシギバナ", "おだやか", "しんりょく", "♀", "くろいヘドロ",
		MoveNames{"ギガドレイン", "ヘドロばくだん", "やどりぎのタネ", "まもる"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: MAX_EFFORT_VAL, SpDef: MAX_EFFORT_VAL, Speed: 4})

	if err != nil {
		panic(err)
	}
	return pokemon
}

//6
func NewCharizard() Pokemon {
	pokemon, err := NewPokemon("リザードン", "ひかえめ", "もうか", "♂", "こだわりスカーフ",
		MoveNames{"オーバーヒート", "エアスラッシュ", "りゅうのはどう", "かえんほうしゃ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 4, SpAtk: MAX_EFFORT_VAL, Speed: MAX_EFFORT_VAL})

	if err != nil {
		panic(err)
	}
	return pokemon
}

//9
func NewBlastoise() Pokemon {
	pokemon, err := NewPokemon("カメックス", "ひかえめ", "げきりゅう", "♂", "しろいハーブ",
		MoveNames{"なみのり", "れいとうビーム", "あくのはどう", "からにこもる"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 4, SpAtk: MAX_EFFORT_VAL, Speed: MAX_EFFORT_VAL})

	if err != nil {
		panic(err)
	}
	return pokemon
}

//91
func NewCloyster() Pokemon {
	pokemon, err := NewPokemon("パルシェン", "やんちゃ", "スキルリンク", "♂", "いのちのたま",
		MoveNames{"つららおとし", "からをやぶる", "ロックブラスト", "ハイドロポンプ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP:4, Atk: MAX_EFFORT_VAL, Speed: MAX_EFFORT_VAL},
	)

	if err != nil {
		panic(err)
	}
	return pokemon
}

//94
func NewGengar() Pokemon {
	pokemon, err := NewPokemon("ゲンガー", "おくびょう", "のろわれボディ", "♀", "くろいヘドロ",
		MoveNames{"たたりめ", "おにび", "みがわり", "かなしばり"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: MAX_EFFORT_VAL, SpAtk:4, Speed: MAX_EFFORT_VAL},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//130
func NewGyarados() Pokemon {
	pokemon, err := NewPokemon("ギャラドス", "いじっぱり", "いかく", "♀", "オボンのみ",
		MoveNames{"たきのぼり", "こおりのキバ", "りゅうのまい", "ちょうはつ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP:4, Atk: MAX_EFFORT_VAL, Speed: MAX_EFFORT_VAL},
	)

	if err != nil {
		panic(err)
	}
	return pokemon
}

//145
func NewZapdos() Pokemon {
	pokemon, err := NewPokemon("サンダー", "ひかえめ", "プレッシャー", "不明", "たべのこし",
		MoveNames{"10まんボルト", "ぼうふう", "ねっぷう", "はねやすめ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP:252, SpAtk:252, Speed:4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//212
func NewScizor() Pokemon {
	pokemon, err := NewPokemon("ハッサム", "いじっぱり", "テクニシャン", "♂", "こだわりハチマキ",
		MoveNames{"とんぼがえり", "バレットパンチ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: MAX_EFFORT_VAL, Atk: MAX_EFFORT_VAL, Speed: 4},
	)

	if err != nil {
		panic(err)
	}
	return pokemon
}

//214
func NewHeracross() Pokemon {
	pokemon, err := NewPokemon("ヘラクロス", "いじっぱり", "こんじょう", "♂", "こだわりハチマキ",
		MoveNames{"インファイト", "メガホーン", "ストーンエッジ", "シャドークロー"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: MAX_EFFORT_VAL, Atk: MAX_EFFORT_VAL, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//227
func NewSkarmory() Pokemon {
	pokemon, err := NewPokemon("エアームド", "わんぱく", "がんじょう", "♀", "ゴツゴツメット",
		MoveNames{"ボディプレス", "はねやすめ", "ステルスロック", "てっぺき"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: MAX_EFFORT_VAL, Def: MAX_EFFORT_VAL, Speed: 4},
	)

	if err != nil {
		panic(err)
	}
	return pokemon
}

//242
func NewBlissey() Pokemon {
	pokemon, err := NewPokemon("ハピナス", "ひかえめ", "てんのめぐみ", "♀", "たべのこし",
		MoveNames{"トライアタック", "シャドーボール", "かえんほうしゃ", "タマゴうみ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{SpAtk: MAX_EFFORT_VAL, SpDef: MAX_EFFORT_VAL, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//248
func NewTyranitar() Pokemon {
	pokemon, err := NewPokemon("バンギラス", "いじっぱり", "すなおこし", "♀", "こだわりハチマキ",
		MoveNames{"かみくだく", "いわなだれ", "ばかぢから", "アイアンヘッド"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: MAX_EFFORT_VAL, Atk: MAX_EFFORT_VAL, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//381
func NewLatios() Pokemon {
	pokemon, err := NewPokemon("ラティオス", "ひかえめ", "ふゆう", "♂", "こだわりメガネ",
		MoveNames{"りゅうせいぐん", "マジカルフレイム", "10まんボルト", "サイコキネシス"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{SpAtk: MAX_EFFORT_VAL, Speed: MAX_EFFORT_VAL, Def: 4},
	)

	if err != nil {
		panic(err)
	}
	return pokemon
}

//445
func NewGarchomp() Pokemon {
	pokemon, err := NewPokemon("ガブリアス", "ようき", "さめはだ", "♂", "きあいのタスキ",
		MoveNames{"じしん", "スケイルショット", "がんせきふうじ", "ステルスロック"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{Atk: MAX_EFFORT_VAL, Speed: MAX_EFFORT_VAL, HP: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

//488
func NewCresselia() Pokemon {
	pokemon, err := NewPokemon("クレセリア", "ずぶとい", "ふゆう", "♀", "ゴツゴツメット",
		MoveNames{"サイコキネシス", "れいとうビーム", "ムーンフォース", "つきのひかり"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: MAX_EFFORT_VAL, Def: MAX_EFFORT_VAL, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

var NEW_RENTAL_POKEMONS = map[PokeName]func() Pokemon {
	"フシギバナ":NewVenusaur,
	"リザードン":NewCharizard,
	"カメックス":NewBlastoise,
	"パルシェン":NewCloyster,
	"ゲンガー":NewGengar,
	"ギャラドス":NewGyarados,
	"サンダー":NewZapdos,
	"ハッサム":NewScizor,
	"ヘラクロス":NewHeracross,
	"エアームド":NewSkarmory,
	"ハピナス":NewBlissey,
	"バンギラス":NewTyranitar,
	"ラティオス":NewLatios,
	"ガブリアス":NewGarchomp,
	"クレセリア":NewCresselia,
}

func NewRentalTeam(pokeNames ...PokeName) (Team, error) {
	pokemons := make(Team, len(pokeNames))
	for i, pokeName := range pokeNames {
		pokemons[i] = NEW_RENTAL_POKEMONS[pokeName]()
	}
	return NewTeam(pokemons)
}

func NewRentalFighters(pokeName1, pokeName2, pokeName3 PokeName) Fighters {
	return Fighters{NEW_RENTAL_POKEMONS[pokeName1](), NEW_RENTAL_POKEMONS[pokeName2](), NEW_RENTAL_POKEMONS[pokeName3]()}
}

func init() {
	for pokeName, _ := range NEW_RENTAL_POKEMONS {
		if _, ok := POKEDEX[pokeName]; !ok {
			errMsg := fmt.Sprintf("RENTAL_POKEMONS の Key に 存在する %v は たぶんタイピングミスってるんで、作者に連絡よろろん twitter @chuusotunamapai", pokeName)
			fmt.Println(errMsg)
		}
	}
}
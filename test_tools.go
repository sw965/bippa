package bippa

func NewTestVenusaur() Pokemon {
  result, err := NewPokemon(
    "フシギバナ", "しんちょう", "しんりょく", "♀", "くろいヘドロ",
    MoveNames{"ギガドレイン", "やどりぎのタネ", "まもる", "ヘドロばくだん"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &HD252_S4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestCharizard() Pokemon {
  result, err := NewPokemon(
    "リザードン", "おくびょう", "もうか", "♂", "こだわりスカーフ",
    MoveNames{"かえんほうしゃ", "エアスラッシュ", "りゅうのはどう", "オーバーヒート"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &CS252_H4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestBlastoise() Pokemon {
  result, err := NewPokemon(
    "カメックス", "ひかえめ", "げきりゅう",
    "♂", "しろいハーブ", MoveNames{"ハイドロカノン", "からをやぶる", "れいとうビーム", "あくのはどう"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &CS252_H4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestAerodactyl() Pokemon {
  result, err := NewPokemon(
    "プテラ", "ようき", "プレッシャー", "♂", "きあいのタスキ",
    MoveNames{"がんせきふうじ", "じしん", "ステルスロック", "ちょうはつ"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &AS252_H4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestNinetales() Pokemon {
  result, err := NewPokemon(
    "キュウコン", "おくびょう", "ひでり", "♀", "きあいのタスキ",
    MoveNames{"マジカルフレイム", "おにび", "ソーラービーム", "おきみやげ"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &CS252_H4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestBlissey() Pokemon {
  result, err := NewPokemon(
    "ハピナス", "ひかえめ", "てんのめぐみ", "♀", "たべのこし",
    MoveNames{"トライアタック", "シャドーボール", "タマゴうみ", "れいとうビーム"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &CB252_S4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestTyranitar() Pokemon {
  result, err := NewPokemon(
    "バンギラス", "いじっぱり", "すなおこし", "♂", "とつげきチョッキ",
    MoveNames{"ストーンエッジ", "かみくだく", "アイアンヘッド", "ばかぢから"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &HA252_S4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestLudicolo() Pokemon {
  result, err := NewPokemon(
    "ルンパッパ", "おだやか", "あめうけざら", "♀", "たべのこし",
    MoveNames{"なみのり", "あまごい", "やどりぎのタネ", "まもる"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &HD252_S4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestToxicroak() Pokemon {
  result, err := NewPokemon(
    "ドクロッグ", "ようき", "かんそうはだ", "♀", "きあいのタスキ",
    MoveNames{"ヘドロばくだん", "こごえるかぜ", "ちょうはつ", "どくどく"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &CS252_H4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestCresselia() Pokemon {
  result, err := NewPokemon(
    "クレセリア", "ずぶとい", "ふゆう", "♀", "オボンのみ",
    MoveNames{"サイコキネシス", "れいとうビーム", "つきのひかり", "みかづきのまい"},
    ALL_MAX_POINT_UPS[MAX_MOVESET_LENGTH], &ALL_MAX_INDIVIDUAL, &HB252_S4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

var TEST_POKEMONS = map[PokeName]func()Pokemon{
	"フシギバナ":NewTestVenusaur,
	"リザードン":NewTestCharizard,
	"カメックス":NewTestBlastoise,
  "キュウコン":NewTestNinetales,
	"プテラ":NewTestAerodactyl,
  "バンギラス":NewTestTyranitar,
  "ルンパッパ":NewTestLudicolo,
  "ドクロッグ":NewTestToxicroak,
  "クレセリア":NewTestCresselia,
}

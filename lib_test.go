package valpass_test

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/tlinden/valpass"
)

type Tests struct {
	name string
	want bool
	opts valpass.Options
}

var pass_random_good = []string{
	`5W@'"5b5=S)b]):xwBuEEu=,x}A46<aS`,
	`QAfwWn;]6ECn-(wZ-z7MxZL)zRA!TO%t`,
	`_5>}+RMm=FRj1a>r/!gG*3tQ>s<&Uh{I`,
	`~Dc6RHW?Yj"nDj)WaWAg#F<IsA[4j?G{`,
	`B;S0|lq:Ns#!{r1UaE0QG7R}tA'K'TNW`,
	`/~]-bT':EeA:dK&[+752EKvS@C1\U70d`,
	`3>cNh2_1(gB(DsA]m$4f/[hHf>{}E*\Q`,
	`Gr5#qF/!:ih?n7p|c?pN50IWc]5$+Q(]`,
	`S#(|irk.%U}[RBFZ2L;}XdDrmOU;SP<\`,
	`+L:T#&@ce[yqWZ0mTfm[D'#a=Ke[j7w'`,
	`:N8vqQ{Vb]@.y?\P2d8,)yHHE?>l|Gi_`,
	`^+s5,#2h<,?_s_Qsd2l;|D42TV3h{7M^`,
	`.^e#(l5$3}1l/-/Uk0,;t^Z[$X0,'h)O`,
	`]-xAyz-"P$98_Z[77@bmo9ZF)I#"Fa,6`,
	`HLkM\]n70U2qU)%Mp{gK@CHt,twiPzH%`,
	`wU2?2&4yx/7HuR@k:~]%/77,DyaNW|"Q`,
	`nb\ZmKT[J)%@=\nF9E2!%N-(+S}Lq95B`,
	`=+0b2[#FMcT~re:PifIWh$IL+>4uyBg1`,
	`xEm]AS#<]cgayw)>O/c<i,)BO[MC0qF,`,
	`EScP'NqM|7/>7e2'orRcS%x6v[sgX(!p`,
	`[.L|hvRRd.@)y?dH?Z46EcEa%/#!m39j`,
	`,$88R.N+C>+adUcw!D"11$H">:SKOiKp`,
	`8#uY]ByJ]iCNp?6-#;&m\pO[G>*!27ge`,
	`@UNu)/qMT{ekO(}qhh4!HI9\QRdrdh^'`,
	`FfoO3pLr_aoGC]lpvo"?RT3E@2f8-764`,
	`Us.dn65ZmF]M}e0Z!$!r0ex-/Z5nwx?J`,
	`e6p{,373[@c@/:CcQ"+(u^U"}^CzxRY.`,
	`kwpHHIqcsuWOio@jlIA2UbO63dkhh'|D`,
	`Yeq@?/Fq.}}"i2dXT=vR2C56hY9R)!_w`,
	`49ZFp54$@\kJ:D;[ZV(VcY|!&sI\O8;&`,
	`SK(ILi(q#FD-*uBbX4,;;1MM2</Md57(`,
	`TA"s$ix&5tlHqk^)182870PpW4X8jH_]`,
	`i"0&lJa?FA>]sD#:AVI)O7|L2x$$WI>(`,
	`_ao{jJ4Z0#njg}GCV{UpQsQubgb!F$-?`,
	`KtAkA~]c}0gj)H7.C6is*>50eIT$OW?*`,
	`Cv[o<mOux790E|[kNrh<n;S\1qU42kNN`,
	`Xj3:.j%kN?k_qYkNMUcQJe@[<K6v.4R~`,
	`aNRU-vO~LX~AwFbUe9t}[WK*3r;PGc/b`,
	`|E|Jl]YjM<4gNh0b1%)^SP:_;%#A\b4b`,
	`Q4#U1/2'5V[_CzYdm7OSZJJE-cSf9^cG`,
	`!jK6zb4)pGrAL/|w|#$a}O||C(0:>:.6`,
	`7t&/B36m8IeM*^e}.)-/X+M8r7'\q:cu`,
	`0iw8o:,bQJ=;d&<CK6?UcaqggQ&r!~%E`,
	`^/FPWoYDwij"B//t}|3aV6vaLI$\3E4%`,
	`^sJ~J.>r?$u'0J,2VD6$Fou,[D~q_vzO`,
	`rVV\wI.L@AAI?+;lU@gnmxKFiob>?s!8`,
	`o]K;x.6$u|^M7kL:lM"13a@rQiD1IJoh`,
	`xM;!)\?;=!lH]|j^jzGG}?6v*O:s~*o=`,
	`f"7#AnRu*b9_=sk^^mMX?+K^ElemvJ(<`,
	`L4WSx8ocC1$74A4#zF!*h8Bq_Eq/1s7s`,
}

var pass_diceware_good = []string{
	`abutting Eucharist dramatized unlearns`,
	`Terrence decorates dwarfed saucing`,
	`swamping nauseated tapioca ascribe`,
	`insatiably ensconcing royally Clarice`,
	`inshore watchdog blunderers methods`,
	`Plasticine brotherly prances dryness`,
	`rustproof flipper commodity nudging`,
	`unburdened frostings adapter vivider`,
	`facile Niamey begrudge menage`,
	`nightcaps miniseries Hannibal strongly`,
	`foresails produces sufficing cannibal`,
	`berths allowing Lewiston sounds`,
	`hazier Hockney snobbier redefines`,
	`Monroe castaways narwhals roadbed`,
	`schuss Trieste assist kebobs`,
	`anteater pianos damping attaining`,
	`desisting colossus refused Madagascan`,
	`misguiding urinalyses moonscapes Taiping`,
	`fracases Indies dishwasher crimsons`,
	`doorman Kleenexes hostessed stooped`,
	`telephoto boozing monoxide Asiago`,
	`completed dogfish rawboned curvacious`,
	`physics virtually rocketing relevant`,
	`infantile sharpest buckler gazillions`,
	`forbids midlands accosts furniture`,
	`concocts Alcestis nitpicker Hindustan`,
	`heirlooms wending Borodin billows`,
	`commotion absinthe chilis drainer`,
	`prerecord brokerages colonel implied`,
	`spoons abates swathed Pocono`,
	`speedy poultices Smollett tracing`,
	`viragoes unwind gasped earache`,
	`rulings Mencken damasking matched`,
	`Sarajevo footbridge stables furloughed`,
	`proclaimed baffling carefully Anatolia`,
	`Cecily Nicaraguan excrete lobbed`,
	`enfold cranny tearjerker blazon`,
	`bucketed Corneille eclectic Maurine`,
	`Berwick gasohol slices bonkers`,
	`swearers iodized Ohioans warden`,
	`Cortez insular several phloem`,
	`assented insolvent beguile aquaplane`,
	`commend trails Amazon clambering`,
	`excretory greatness plackets creeks`,
	`transistor exclusion inboxes sidling`,
	`cherries elongating Lollard piques`,
	`heartening orbiting zombie revile`,
	`reconcile completes roughs innocence`,
	`quickness Cheever Thimbu scours`,
	`hobble piteously precepts sorest`,
	`braving shirted backstage Taiping`,
}

var pass_worst_bad = []string{
	`123456`, `charlie`, `summer`, `sophie`, `merlin`,
	`password`, `aa123456`, `George`, `Ferrari`, `cookie`,
	`123456789`, `donald`, `Harley`, `Cheese`, `ashley`,
	`12345678`, `password1`, `222222`, `Computer`, `bandit`,
	`12345`, `qwerty123`, `Jessica`, `jesus`, `killer`,
	`111111`, `letmein`, `ginger`, `Corvette`, `aaaaaa`,
	`1234567`, `zxcvbnm`, `abcdef`, `Mercedes`, `1q2w3e`,
	`sunshine`, `login`, `Jordan`, `flower`, `zaq1zaq1`,
	`qwerty`, `starwars`, `55555`, `Blahblah`, `mustang`,
	`iloveyou`, `121212`, `Tigger`, `Maverick`, `test`,
	`princess`, `bailey`, `Joshua`, `Hello`, `hockey`,
	`admin`, `freedom`, `Pepper`, `loveme`, `dallas`,
	`welcome`, `shadow`, `Robert`, `nicole`, `whatever`,
	`666666`, `passw0rd`, `Matthew`, `hunter`, `admin123`,
	`abc123`, `master`, `12341234`, `amanda`, `michael`,
	`football`, `baseball`, `Andrew`, `jennifer`, `liverpool`,
	`123123`, `buster`, `lakers`, `banana`, `querty`,
	`monkey`, `Daniel`, `andrea`, `chelsea`, `william`,
	`654321`, `Hannah`, `1qaz2wsx`, `ranger`, `soccer`,
	`!@#$%^&*`, `Thomas`, `starwars`, `trustno1`, `london`,
}

var pass_dict_bad = []string{
	`clued`, `lads`, `stifle`,
	`receptivity`, `apprehends`, `accounts`,
	`putts`, `spurt`, `sideswipe`,
	`dabbed`, `goatskin`, `nooks`,
	`sulkiness`, `worships`, `coevals`,
	`entwining`, `sportscasters`, `pew`,
	`horse`, `daybeds`, `booklet`,
	`Suzette`, `abbreviate`, `stubborn`,
	`govern`, `ageism`, `refereeing`,
	`dents`, `Wyeth`, `concentric`,
	`Kamehameha`, `grosser`, `belie`,
	`wherefore`, `president`, `pipit`,
	`pinholes`, `mummifying`, `quartermasters`,
	`fruitlessness`, `seafarer`, `Einsteins`,
	`stomping`, `glided`, `retried`,
	`effected`, `ministry`,
}

var opts_std = valpass.Options{
	Compress:         valpass.MIN_COMPRESS,
	CharDistribution: valpass.MIN_DIST,
	Entropy:          valpass.MIN_ENTROPY,
	Dictionary:       nil,
	UTF8:             false,
}

var opts_dict = valpass.Options{
	Compress:         valpass.MIN_COMPRESS,
	CharDistribution: valpass.MIN_DIST,
	Entropy:          valpass.MIN_ENTROPY,
	Dictionary:       &valpass.Dictionary{Words: ReadDict("t/american-english")},
	UTF8:             false,
}

var goodtests = []Tests{
	{
		name: "checkgood",
		want: true,
		opts: opts_std,
	},
	{
		name: "checkgood-dict",
		want: true,
		opts: opts_dict,
	},
}

var badtests = []Tests{
	{
		name: "checkbad",
		want: false,
		opts: opts_std,
	},
	{
		name: "checkbad-dict",
		want: false,
		opts: opts_dict,
	},
}

func TestValidate(t *testing.T) {
	t.Parallel()

	for _, tt := range goodtests {
		for _, pass := range pass_random_good {
			CheckPassword(t, pass, tt.name, tt.want, tt.opts)
		}

		for _, pass := range pass_diceware_good {
			CheckPassword(t, pass, tt.name, tt.want, tt.opts)
		}
	}

	for _, tt := range badtests {
		for _, pass := range pass_worst_bad {
			CheckPassword(t, pass, tt.name, tt.want, tt.opts)
		}

		for _, pass := range pass_dict_bad {
			CheckPassword(t, pass, tt.name, tt.want, tt.opts)
		}
	}

}

func CheckPassword(t *testing.T, password string,
	name string, want bool, opts valpass.Options) {

	result, err := valpass.Validate(password, opts)
	if err != nil {
		t.Errorf("test %s failed with error: %s\n", name, err)
	}

	if want && !result.Ok {
		t.Errorf("test %s failed. pass: %s, want: %t, got: %t, dict: %t\nresult: %v\n",
			name, password, want, result.Ok, result.DictionaryMatch, result)
	}

	if !want && result.Ok {
		t.Errorf("test %s failed. pass: %s, want: %t, got: %t, dict: %t\nresult: %v\n",
			name, password, want, result.Ok, result.DictionaryMatch, result)
	}
}

func BenchmarkValidateEntropy(b *testing.B) {
	passwords := GetPasswords(b.N)

	for i := 0; i < b.N; i++ {
		_, err := valpass.Validate(passwords[i], valpass.Options{Entropy: 10})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkValidateCharDist(b *testing.B) {
	passwords := GetPasswords(b.N)

	for i := 0; i < b.N; i++ {
		_, err := valpass.Validate(passwords[i], valpass.Options{CharDistribution: 10})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkValidateCompress(b *testing.B) {
	passwords := GetPasswords(b.N)

	for i := 0; i < b.N; i++ {
		_, err := valpass.Validate(passwords[i], valpass.Options{Compress: 10})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkValidateDict(b *testing.B) {
	passwords := GetPasswords(b.N)

	for i := 0; i < b.N; i++ {
		_, err := valpass.Validate(passwords[i],
			valpass.Options{Dictionary: &valpass.Dictionary{Words: ReadDict("t/american-english")}},
		)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkValidateAll(b *testing.B) {
	passwords := GetPasswords(b.N)

	for i := 0; i < b.N; i++ {
		_, err := valpass.Validate(passwords[i])
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkValidateAllwDict(b *testing.B) {
	passwords := GetPasswords(b.N)

	for i := 0; i < b.N; i++ {
		_, err := valpass.Validate(passwords[i], opts_dict)
		if err != nil {
			panic(err)
		}
	}
}

func ReadDict(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func GetPasswords(count int) []string {

	cmd := exec.Command("pwgen", "-1", "-s", "-y", "32", fmt.Sprintf("%d", count+1))

	out, err := cmd.Output()
	if err != nil {
		panic(cmd.Err)
	}

	return strings.Split(string(out), "\n")
}

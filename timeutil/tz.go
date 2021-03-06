package timeutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/tkuchiki/go-timezone"
)

var (
	Timezone *timezone.Timezone = timezone.New()

	timezones = map[string]bool{
		"ACDT":   true,
		"ACST":   true,
		"ACT":    true,
		"ACWDT":  true,
		"ACWST":  true,
		"ACWT":   true,
		"ADT":    true,
		"AEDT":   true,
		"AEST":   true,
		"AET":    true,
		"AFT":    true,
		"AKDT":   true,
		"AKST":   true,
		"AKT":    true,
		"ALMST":  true,
		"ALMT":   true,
		"AMST":   true,
		"AMT":    true,
		"ANAT":   true,
		"AQTST":  true,
		"AQTT":   true,
		"ARST":   true,
		"ART":    true,
		"AST":    true,
		"AWDT":   true,
		"AWST":   true,
		"AWT":    true,
		"AZOST":  true,
		"AZOT":   true,
		"AZST":   true,
		"AZT":    true,
		"BDST":   true,
		"BDT":    true,
		"BNT":    true,
		"BORTST": true,
		"BOT":    true,
		"BRST":   true,
		"BRT":    true,
		"BST":    true,
		"BTT":    true,
		"CAST":   true,
		"CAT":    true,
		"CCT":    true,
		"CDT":    true,
		"CEST":   true,
		"CET":    true,
		"CHADT":  true,
		"CHAST":  true,
		"CHOST":  true,
		"CHOT":   true,
		"CHST":   true,
		"CHUT":   true,
		"CKHST":  true,
		"CKT":    true,
		"CLST":   true,
		"CLT":    true,
		"COST":   true,
		"COT":    true,
		"CST":    true,
		"CT":     true,
		"CVT":    true,
		"CXT":    true,
		"DAVT":   true,
		"DDUT":   true,
		"EASST":  true,
		"EAST":   true,
		"EAT":    true,
		"ECT":    true,
		"EDT":    true,
		"EEST":   true,
		"EET":    true,
		"EGST":   true,
		"EGT":    true,
		"EHDT":   true,
		"EST":    true,
		"FJST":   true,
		"FJT":    true,
		"FKT":    true,
		"FNST":   true,
		"FNT":    true,
		"GALT":   true,
		"GAMT":   true,
		"GDT":    true,
		"GET":    true,
		"GFT":    true,
		"GHST":   true,
		"GILT":   true,
		"GMT-1":  true,
		"GMT-10": true,
		"GMT-11": true,
		"GMT-12": true,
		"GMT-2":  true,
		"GMT-3":  true,
		"GMT-4":  true,
		"GMT-5":  true,
		"GMT-6":  true,
		"GMT-7":  true,
		"GMT-8":  true,
		"GMT-9":  true,
		"GMT":    true,
		"GMT+1":  true,
		"GMT+10": true,
		"GMT+11": true,
		"GMT+12": true,
		"GMT+13": true,
		"GMT+14": true,
		"GMT+2":  true,
		"GMT+3":  true,
		"GMT+4":  true,
		"GMT+5":  true,
		"GMT+6":  true,
		"GMT+7":  true,
		"GMT+8":  true,
		"GMT+9":  true,
		"GST":    true,
		"GYT":    true,
		"HADT":   true,
		"HAST":   true,
		"HAT":    true,
		"HKST":   true,
		"HKT":    true,
		"HOVST":  true,
		"HOVT":   true,
		"ICT":    true,
		"IDT":    true,
		"IOT":    true,
		"IRDT":   true,
		"IRKST":  true,
		"IRKT":   true,
		"IRST":   true,
		"IRT":    true,
		"IST":    true,
		"JDT":    true,
		"JST":    true,
		"KDT":    true,
		"KGT":    true,
		"KOST":   true,
		"KRAST":  true,
		"KRAT":   true,
		"KST":    true,
		"LHDT":   true,
		"LHST":   true,
		"LHT":    true,
		"LINT":   true,
		"MAGST":  true,
		"MAGT":   true,
		"MALST":  true,
		"MART":   true,
		"MAWT":   true,
		"MDT":    true,
		"MEST":   true,
		"MET":    true,
		"MHT":    true,
		"MIST":   true,
		"MLAST":  true,
		"MMT":    true,
		"MSD":    true,
		"MSK":    true,
		"MST":    true,
		"MUST":   true,
		"MUT":    true,
		"MVT":    true,
		"MYT":    true,
		"NCST":   true,
		"NCT":    true,
		"NDT":    true,
		"NFDT":   true,
		"NFT":    true,
		"NOVT":   true,
		"NPT":    true,
		"NRT":    true,
		"NST":    true,
		"NT":     true,
		"NUT":    true,
		"NZDT":   true,
		"NZST":   true,
		"NZT":    true,
		"OMSST":  true,
		"OMST":   true,
		"ORAT":   true,
		"PDT":    true,
		"PEST":   true,
		"PET":    true,
		"PETT":   true,
		"PGT":    true,
		"PHOT":   true,
		"PHST":   true,
		"PHT":    true,
		"PKST":   true,
		"PKT":    true,
		"PMDT":   true,
		"PMST":   true,
		"PONT":   true,
		"PST":    true,
		"PT":     true,
		"PWT":    true,
		"PYST":   true,
		"PYT":    true,
		"QYZST":  true,
		"QYZT":   true,
		"RET":    true,
		"ROTT":   true,
		"SAKT":   true,
		"SAMT":   true,
		"SAST":   true,
		"SBT":    true,
		"SCT":    true,
		"SGT":    true,
		"SRET":   true,
		"SRT":    true,
		"SST":    true,
		"SYOT":   true,
		"TAHT":   true,
		"TFT":    true,
		"TJT":    true,
		"TKT":    true,
		"TLT":    true,
		"TMT":    true,
		"TOST":   true,
		"TOT":    true,
		"TRT":    true,
		"TSD":    true,
		"TVT":    true,
		"ULAST":  true,
		"ULAT":   true,
		"UTC":    true,
		"UYST":   true,
		"UYT":    true,
		"UZST":   true,
		"UZT":    true,
		"VET":    true,
		"VLAST":  true,
		"VLAT":   true,
		"VOLT":   true,
		"VOST":   true,
		"VUST":   true,
		"VUT":    true,
		"WAKT":   true,
		"WAT":    true,
		"WEST":   true,
		"WET":    true,
		"WFT":    true,
		"WGST":   true,
		"WGT":    true,
		"WIB":    true,
		"WIT":    true,
		"WITA":   true,
		"WSDT":   true,
		"WST":    true,
		"YAKST":  true,
		"YAKT":   true,
		"YEKST":  true,
		"YEKT":   true,
	}
)

func GetTimezoneLocation(s string) (*time.Location, error) {
	code := strings.ToUpper(s)
	if _, ok := timezones[code]; !ok {
		return nil, fmt.Errorf("invalid timezone %q", s)
	}
	names := Timezone.Timezones()[code]
	if len(names) == 0 {
		return nil, fmt.Errorf("no location data for %q", code)
	}
	return time.LoadLocation(names[0])
}

func ReTimezoneCodes() string {
	codes := make([]string, 0, len(timezones))
	for k := range timezones {
		codes = append(codes, strings.ReplaceAll(k, "+", "\\+"))
	}
	return strings.Join(codes, "|")
}

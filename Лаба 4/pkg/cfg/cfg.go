package cfg

import "time"

// files
const (
	LinksFileName   = "Links.xlsx"
	ResultExcelFile = "result.xlsx"
	ResultCSVFile   = "result.csv"
)

// excel
const (
	SheetName = "Лист1"
)

// vessel parsing
const (
	RequestDelay = 500 * time.Millisecond
	HTTPTimeout  = 20 * time.Second

	BaseURL         = "https://www.vesselfinder.com"
	VesselPath      = "/ru/vessels/details/"
	ShipLinkClass   = "a.ship-link"
	TitleSelector   = "h1.title"
	IMOMMSISelector = "td.v3.v3np"
	TypeSelector    = "td.n3"
	AISTypeText     = "AIS тип"
)

package collection

type Collection struct {
	Info     *Schema     `json:"info"`
	Item     *[]*Node    `json:"item"`
	Auth     *Auth       `json:"auth"`
	Variable *[]Variable `json:"variable"`
}

type Schema struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type Node struct {
	Name     string    `json:"name"`
	Item     *[]*Node  `json:"item"`
	Request  *Request  `json:"request"`
	Response *[]string `json:"response"`
}

type Request struct {
	Method string   `json:"method"`
	Header []string `json:"header"`
	Url    URL      `json:"url"`
	Body   *Body    `json:"body"`
	Auth   *Auth    `json:"auth"`
}

type URL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type Body struct {
	Mode    string  `json:"mode"`
	Raw     string  `json:"raw"`
	Options Options `json:"options"`
}

type Options struct {
	Raw Raw `json:"raw"`
}

type Raw struct {
	Language string `json:"language"`
}

type Auth struct {
	Type   string  `json:"type"`
	Bearer *Bearer `json:"bearer"`
}

type Bearer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

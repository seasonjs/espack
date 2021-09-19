package config

// EntryPoints https://webpack.docschina.org/concepts/entry-points/
type EntryPoints interface{}

// OutputPoints  https://webpack.docschina.org/concepts/output/#root
type OutputPoints interface{}

// PluginPoints https://webpack.docschina.org/concepts/plugins/#anatomy

type PluginPoints interface{}

type configuration struct {
	// string
	Entry EntryPoints `json:"entry"`
	//string []string
	Output OutputPoints `json:"output"`

	Plugin PluginPoints `json:"plugin"`
}

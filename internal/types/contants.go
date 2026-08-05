package types

var ExtensionOverrides = map[string]string{
	"html":            "HTML",
	"ecmarkup":        "HTML",
	"xml":             "XML",
	"yaml":            "YAML",
	"yml":             "YAML",
	"miniyaml":        "YAML",
	"gerber image":    "Solidity",
	"sol":             "Solidity",
	"renderscript":    "Rust",
	"rs":              "Rust",
	"properties":      "Java Properties",
	"ini":             "Java Properties",
	"java-properties": "Java Properties",
	"db":              "Database",

	//  THE JAVASCRIPT / TYPESCRIPT ECOSYSTEM
	"js":               "JavaScript",
	"javascript":       "JavaScript",
	"javascriptreact":  "JavaScript React",
	"jsx":              "JavaScript React",
	"javascript react": "JavaScript React", // Prevents re-casing bugs!

	"ts":               "TypeScript",
	"typescript":       "TypeScript",
	"typescriptreact":  "TypeScript React",
	"tsx":              "TypeScript React",
	"typescript react": "TypeScript React", // Prevents re-casing bugs!

	// SHELL & TERMINAL
	"sh":          "Shell Script",
	"shellscript": "Shell Script",
	"shell":       "Shell Script",
	"zsh":         "Shell Script",

	//  CONFIGS & ENV
	"conf":          "Configuration",
	"configuration": "Configuration",
	"dotenv":        "Environment File",
	".env":          "Environment File",

	// GIT & WORKFLOWS
	"github-actions-workflow": "GitHub Actions",
	"github actions":          "GitHub Actions",
	"ignore":                  "Git Config",
	"ignore list":             "Git Config",
	"gitignore file":          "Git Config",

	// PYTHON DATA
	"jupyter-notebook": "Jupyter Notebook",
	"jupyter notebook": "Jupyter Notebook",
	"pip-requirements": "Pip Requirements",
	"pip requirements": "Pip Requirements",

	//  GO DATA
	"go":        "Go",
	"mod":       "Go",
	"go.mod":    "Go",
	"go module": "Go",

	//c and c++ family
	"m":   "Objective-C",
	"c":   "C",
	"cpp": "C++",
	"h":   "C",
	"hpp": "C++",

	// PLAIN TEXT VARIATIONS
	"text":       "Plain Text",
	"txt":        "Plain Text",
	"plaintext":  "Plain Text",
	"plain text": "Plain Text",
	"":           "Plain Text",

	//Markdown
	"md": "Markdown",
}

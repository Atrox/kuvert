package kuvert

// Email contains all the information needed to generate an email.
type Email struct {
	Summary    string   // Email summary that will be displayed as a preview in most email clients
	Title      string   // Title of the email
	Intros     []string // Intros to be displayed at the beginning of the email
	Dictionary []Entry  // Dictionary to be displayed as a key-value list
	Table      Table    // Table to be displayed
	Actions    []Action // Actions to be displayed as buttons
	Outros     []string // Outros to be displayed at the end of the email
	Signature  string   // Signature to be displayed at the end of the email
}

func (e *Email) fillDefaults() {
	if e.Title == "" {
		e.Title = "Hey,"
	}
	if e.Signature == "" {
		e.Signature = "Yours truly,"
	}
}

// Entry contains a key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Table contains the data to be displayed in a table.
// The data is a 2D slice of Entry. The first dimension is the row and the second dimension is the column.
// The Columns field contains the custom width and alignment of the columns.
//
// Example:
//
//	Table{
//	    Data: [][]Entry{
//	        {{"Name", "John"}, {"Age", "30"}},
//	        {{"Name", "Jane"}, {"Age", "25"}},
//	    },
//	    Columns: Columns{
//	        CustomWidth: map[string]string{
//	            "Name": "100px",
//	            "Age":  "50px",
//	        },
//	        CustomAlignment: map[string]string{
//	            "Name": "left",
//	            "Age":  "right",
//	        },
//	    },
//	}
type Table struct {
	Data    [][]Entry
	Columns Columns
}

// Columns contains the custom width and alignment of the columns.
// The key is the column name and the value is the width or alignment.
type Columns struct {
	CustomWidth     map[string]string
	CustomAlignment map[string]string
}

// Action contains the instructions and the button to be displayed.
type Action struct {
	Instructions string
	Button       Button
}

// Button represents a button in the email.
// The button can be customized with a color and a text color. The color can be any valid CSS color.
type Button struct {
	Text string
	Link string

	Color     string
	TextColor string
}

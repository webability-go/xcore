// XCore for GO v0
// =============================
// 
// The XCore package is used to build basic object for programation, from the original WebAbility Core for PHP program.
// For GO, the actual existing code includes:
//
// - XCache: Application Memory Caches for any purpose,
//
// - XDataset: Basic nested data structures for any purpose (template injection, configuration files, database records, etc),
//
// - XLanguage: language dependant text tables,
//
// - XTemplate: template system with meta language.
package xcore

// VERSION: is the used version nombre of the XCore library.
const VERSION = "0.1.1"

// if LOG is set to TRUE, LOG indicates to the XCore libraries to log a trace of functions called, with most important parameters.
//
// LOG can be set to true or false dynamically to trace only parts of code on demand.
var LOG = false

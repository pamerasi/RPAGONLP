//  ================================================================================================================
//  PROBLEM    : Write a Parser to tokenize/validate Dept, Course, Semester, Year for "Course Selection input text"
//  REQUIRMENT : INPUT "CS111 2016 Fall", "CS-111 Fall 2016", "CS 111 F2016"    OUTPUT: | CS | 111 | Fall | 2016 | 
//  AUTHOR     : Prasanna Amerasinghe 
//  DATE       : 09/01/2019
//  ================================================================================================================
//  Assumptions 
//  (NOTE: [DeptCourse] & [OfferSession] are considerd "Fields". [Dept], [Course], [Year], [Semester] are "Tokens"
//  ----------------------------------------------------------------------------------------------------------------
// 1) Skip leading spaces and Delimiters before the [DeptCourse] Field  (Valid delimiters are ' ' ,  ':'  ,  '-'  )  
// 2) There should be ONE Field Seperator, "a space" between the [DeptCourse] AND [OfferSession] Fields
// 3) There should be either NOTHING or ONE delimiter between [Dept] and [Course] tokens     
// 4) The [OfferSession] Field is either [Year]+[Semester] OR [Semester]+[Year].   Both token orders are supported!
// 5) There could be any number of valid delimiters between [Year] and [Semester] tokens
// 6) [Year] token data is "range validated" to be between 2007 - 2021.    
// 7) [Semester] token data is "lookup validated" using a Dictionary  
//===================================================================================================================
//  Code Outline
//  ------------
// L0 :                                             main() 
//                        |---------------------------^---------------------------|
// L1 :            getDeptCourse()                                        getOfferSession()
//               |---------^-------|                             |----------------^------------------|
// L2 :   getDeptToken()    getCourseToken()             getYearToken()                      getSemesterToken()  
//               |                 |                  |--------^--------|                  |---------^---------|
// L3 :          |                 |                  |            validateYear()          |          validateSemester()
//               |                 |                  |                                    |
// L4 :   getAlphaToken()    getNumberToken()    getNumberToken()                     getAlphaToken()
// -----------------------------------------------------------------------------------------------------
// Utility L-1       :        skipSpacesDelims()   isValid()    isDelim() 
// Utility Primitive :                isNumber()     isLetter()         
//
//===================================================================================================================
// Notes : Focused on a MVP solution using basic golang constructs.     
//       : Detailed error handling to improve customer feedback on data entry errors
//       : funcid based Documentation / Tracing Scheme to improve code maintainability. (To enable set LetsTrace=true)      
//===================================================================================================================
//
//   
package main

import (
	"fmt"
	"strings"
	"strconv"
	"os"
	"bufio"
)

// CH String Type
 type ChStr struct {
 	data string
 	indx int
 	len  int
 }

var InputStrStruct ChStr

// Course Selection Problem-1 Parameters
 const MaxTokens = 10
 const CurTokens = 4 


// Course Selection Problem Token Types
 const Dept      = 0
 const Course    = 1
 const Year      = 2
 const Semester  = 3

// Field Seperator 
// CH-Problem Statement : "There is always "a space" after the Course Number and before Semester+Year
// (Comment : When input data entry is "form based", the "Field Seperator" could ideally be a non-keyboard 
//            character, inserted by mobile/web client code. Can be used to increase parsing concurrency, accuracy. 
//            better error handling etc. e.g. when parsing input string "CS 2018 Fall", Is 2018 a year or a course? 
 var FieldSeperator byte = ' ' 
   

 // Valid range of years for Courses 
 // Assumption : All 2 digit year abbreviations will be normalized to the 21st centuary.  
 // e.g. Abbreviated year ntry 89 will be 2089   
 const EarliestCourseYear = 2007
 const LatestCourseYear   = 2022
 
 
 // Semester Lookup Dictionary
 var ValidSemester = map[string] string {
    "F"         : "Fall",
    "FA"        : "Fall",
    "FAL"       : "Fall",        
    "FALL"      : "Fall",
    "S"         : "Spring",
    "SP"        : "Spring",
    "SPR"       : "Spring",
    "SPRG"      : "Spring", 
    "SPRNG"     : "Spring",           
    "SPRING"    : "Spring",
    "SU"        : "Summer",
    "SUMR"      : "Summer",
    "SUMMER"    : "Summer",
    "W"         : "Winter",
    "WI"        : "Winter",
    "WIN"       : "Winter",  
    "WTR"       : "Winter",    
    "WNTR"      : "Winter", 
    "WINTR"     : "Winter",     
    "WINTER"    : "Winter",
 } 
 
 // Debug, Trace Global Scoped Error handling variables
 // var LetsDebug bool
 var LetsTrace bool
 

//===========================================================
//=== String Parsing & Validation Utility functions =========
//===========================================================
//funcid:100
func isLetter (c byte) bool {
	if ((c >= 'a') && (c <= 'z') || (c >='A') && (c <='Z')) {
		return true
	}
	return false
}

//funcid:130 
func isNumber (c byte) bool {
	if ((c >= '0') && (c <= '9')) {
		return true
	}
	return false
}

//funcid:150
func isDelimiter (c byte) bool {
	if (c == ' ') || (c == '-') || (c == ':') {
		return true
	}
	return false
	
}


//funcid:170
func isValid (c byte) bool {
	if (isDelimiter(c) || isNumber(c) || isLetter(c)) {
		return true
	} else {
		return false
	}
	
}

//===========================================================
//============ Common Token Parsing Functions ===============
//===========================================================

// Skips Spaces Delimiters by advancing "indx" along "data" string
//funcid:500
func skipSpacesDelims(inStr *ChStr) string {
	var char byte
	var err string
	
	if (LetsTrace) {
		fmt.Printf("..TRACE-    500.10 : IN- : skipSpaceDelims() \n")
	}
	
	if (inStr.indx < 0 || inStr.indx >= inStr.len) {
		err = "PANIC-500.20 -  Invalid input structure  ==> '" + string(inStr.indx) + "' " + " \n " + err		
		return err
	}	
	
	for (inStr.indx < inStr.len) {
		char = inStr.data[inStr.indx]
		if !(isValid(char)) {
			err = "ERROR-500.30 - Invalid Character ==> '" +  string(char) + "' " + " \n " + err
			return err
		}
		
		if (isDelimiter(char)) {
			inStr.indx++
		} else {
			break
		}
	} // for
	
	if (LetsTrace) {
		fmt.Printf("..TRACE-    500.90 : OUT : skipSpaceDelims() \n")
	}	
	
	return ""
}



// Parses and Extracts a string of Alpha characters into a return "Token" string
// funcid:600
func getAlphaToken (inStr *ChStr) (string, string) {
	var alphaToken string
	var char byte
	var err string
	
	if (LetsTrace) {
		fmt.Printf("......TRACE-600.10 : IN- : getAlphaToken() \n")
	}
	
	if (inStr.indx < 0 || inStr.indx >= inStr.len) {
		err = "PANIC-600.20 - Invalid input structure  ==> " + string(inStr.indx)				
		return "", err
	}		
	
	if !(isLetter(inStr.data[inStr.indx])) {
		err = "ERROR-600.30 - Non Alpha first character in Alpha Token ==> '" + string(inStr.data[inStr.indx]) + "' " + " \n " + err	
		return "", err
	}
	
	
	for (inStr.indx < inStr.len) {
		char = inStr.data[inStr.indx]
		if !(isValid(char)) {
			err = "ERROR-600.40 - Invalid Character around Alpha token => '" + string(char) + "'" + " \n" + err		
			return "", err
		}
		
		if isLetter(char)  {
			alphaToken += string(char)
			inStr.indx++
			continue
		} else {
			break
		} // if-else				
	} // For Loop - Semester Parser	
	
	if (LetsTrace) {
		fmt.Printf("......TRACE-600.90 : OUT : getAlphaToken() \n")	
	}
	
	return alphaToken, ""

}

// Parses and Extract a string of Numeric characters into a return "Token" string
// funcid:650
func getNumberToken (inStr *ChStr) (string, string) {
	var numberToken string
	var char byte
	var err string
	
	if (LetsTrace) {
		fmt.Printf("......TRACE-650.10 : IN- : getNumberToken() %v \n", inStr)
	}
	
	if (inStr.indx < 0 || inStr.indx >= inStr.len) {
		err = "PANIC-650.20 - Invalid input structure  ==> '" + string(inStr.indx) + "' " + " \n " + err	
		return "", err
	}		
	
	if !(isNumber(inStr.data[inStr.indx])) {
		err = "ERROR-650.30 - Non Number first character in Number Token ==> " + string(inStr.data[inStr.indx]) + "'" + " \n " + err	
		return "", err
	}
	
	
	for (inStr.indx < inStr.len) {
		char = inStr.data[inStr.indx]
		if !(isValid(char)) {
			err = "ERROR-650.40 - Invalid Character around Number token => '" + string(char) + "'" + " \n " + err			
			return "", err
		}
		
		if isNumber(char)  {
			numberToken += string(char)
			inStr.indx++
			continue
		} else {
			break
		} // if-else				
	} // For Loop - Semester Parser	
	
	if (LetsTrace) {
		fmt.Printf("......TRACE-650.90 : OUT : getNumberToken() numberToken %v from inStr %v\n", numberToken, inStr)	
	}
	
	return numberToken, ""

}


//====================================================================
//============ Course Selection Solution specific Functions ==========
//      Functions to parse [DeptCourse] and [OfferSession] Fields
//====================================================================

// Function to parse input string and Process  [DeptCourse] Field and 
// update return the object tokenList

//funcid:700
func getDeptCourse (inStr *ChStr, tokenArr []string ) string {
	var char byte
	var err  string
	 
	
	if (LetsTrace) {
		fmt.Printf("..TRACE-    700.10 : IN- : getDeptCourse() \n")
	}	
	

	// start Processing Department Alpha Token
	if (inStr.indx < 0 || inStr.indx >= inStr.len) {
		err = "PANIC-700.40 - Empty Input or Invalid Input String " + inStr.data	
		return err	
	}
	
	char = inStr.data[inStr.indx]
	
// Invalid Department Format	
	if !(isLetter(char)) {
		err =  "ERROR-700.50 - Department data should have Alpha characters " + inStr.data
		return err					
	}
	
	
 
// Get Department
    if (isLetter(char)) {
		err = getDeptToken(inStr, tokenArr)
		if (err != "") {
			err = "ERROR-700.55 - During or after Parsing Dept "  + " \n " + err 			
		    return err
		}	
    }
					
	if (inStr.indx < 0 || inStr.indx >= inStr.len) {
		err = "ERROR-700.58 - Missing Course Data"     
		return err		
	}	
  
    char = inStr.data[inStr.indx] 
 
 // Skip a single Delimiter, if present    
    if (isDelimiter(char)) {
 	    if (inStr.indx + 1 >= inStr.len) {
			err = "ERROR-700.60 - Missing Course Data in  \n" 
		    return err		
         }	//if 
         
        inStr.indx++ 
        char = inStr.data[inStr.indx]
    } // if (isDelimiter(c))
 
 // Invalid Course Format	
	if !(isNumber(char)) {
		err =  "ERROR-700.63 - Course Entry must start with Numeric characters. Invalid ==> '" + string(char) + "'" + " \n" + err		
		return err					
	}
 
 // Get Course 
   if (isNumber(char)) {
		err = getCourseToken(inStr, tokenArr)
		if (err != "") {
			err = "ERROR-700.65 - After Parsing Course " + " \n " + err    	    
		    return err
		}	//if (err != "")
   } //if (isNumber(char))
	
	if (err != "") {
		err = "ERROR-700.70 - Parsing the [DeptCourse] Field " + " \n " + err	
	}	
	
	
	
	if (LetsTrace) {
		fmt.Printf("..TRACE-    700.90 : OUT : getDeptCourse() \n")
	}
	
	return ""
	

}


// Parse and Extract the Department Token for the [DeptCourse] Field
//funcid:720
func getDeptToken (inStr *ChStr, tokenArr []string ) string {
	 var retToken string
	 var err      string
	 
	 
	 if (LetsTrace) {
		fmt.Printf("....TRACE-  720.10 : IN- : getDeptToken() %v \n", inStr)
	}	
		  
	 
	 retToken, err = getAlphaToken(inStr)
	 if (err != "") {
	 	err = "ERROR-720.30 - When Getting Department data " + " \n " + err	 	
	 	return err
	 }	 
	 
	 tokenArr[Dept] += retToken	
	 
	 
	 
	  if (LetsTrace) {
		fmt.Printf("....TRACE-  720.90 : OUT : getDeptToken() %v retToken %v  tokenArr[Semester]-%v\n", inStr, retToken, tokenArr[Semester])	 
	  }
	 
	 return err

}



// Parse and Extract the Course Token for the [DeptCourse] Field
// funcid:750
func getCourseToken (inStr *ChStr, tokenArr []string ) string {
	 var retToken string
	 var err      string
	 
	 
	 if (LetsTrace) {
		fmt.Printf("....TRACE-  750.10 : IN- : getCourseToken() %v \n", inStr)
	}	
		 
	 
	 retToken, err = getNumberToken(inStr)
	 if (err != "") {
	 	err = "ERROR-750.30 - When Getting Course data " + err 
	 	return err
	 }
	 
	 
	 
	 tokenArr[Course] += retToken	
	
	 	 
	 
	 if (LetsTrace) {
		fmt.Printf("....TRACE-  750.90 : OUT : getCourseToken() %v retToken %v \n", inStr, retToken)
	}	 
	  
	  return err

}

// ========================================================================
// Function to parse input string and Process  [OfferSession] Field and 
// update return the object tokenList
// NOTE: Handles both [Semester-Year] as well as [Year-Semester] formats
// =======================================================================
//funcid:800
func getOfferSession (inStr *ChStr, tokenArr []string ) string {
	var char  byte
	var err   string
	 
	
	if (LetsTrace) {
		fmt.Printf("..TRACE-    800.10 : IN- : getClassSession() \n")
	}	
	
	if (inStr.indx < 0 || inStr.indx >= inStr.len) {
		err =  "PANIC-800.20 - Invalid input structure  ==> " + string(inStr.indx) + " \n " + err
		return err
	}	
	

	char = inStr.data[inStr.indx]


// Invalid Data
	
	if !(isNumber(char) ||  isLetter(char) ) {
		err = "ERROR-800.15 - Found invalid data in getClassSession()  Char '" + string(char) + "'" + " \n " + err
		return err		
	} // if !((isNumber(char) ||  isLetter(char) )
	


// Parse YEAR-SEMESTER Format

	if (isNumber(char)) {
		err = getYearToken(inStr, tokenArr)	
		if (err != "") {
			err = "ERROR-800.25 - When Parsing Year Data " +  " \n " + err  
		    return err
		}
		
		if (inStr.indx < 0 || inStr.indx >= inStr.len) {
			err = "ERROR-800.26 - Missing Semester Data "  + " \n " + err		    
			return err		
		}					
		
		err = skipSpacesDelims(inStr)
		if (err != "") {
			err = "ERROR-800.27 -  Skipping Spaces before Semester " + " \n " + err 	    
		    return err
		}
		
		if (inStr.indx < 0 || inStr.indx >= inStr.len) {
			err = "ERROR-800.28 - Missing Semester Data "  + " \n " + err	    
			return err			
		}
		
		err = getSemesterToken(inStr, tokenArr)
		if (err != "") {
			err = "ERROR-800.29 - in getting Semester  " + " \n " + err 		
		    return err
		}	
	} // if (isNumber(char) //	 
	
	
// Parse SEMESTER-YEAR Format
	
	if (isLetter(char)) {
		err = getSemesterToken(inStr, tokenArr)
		if (err != "") {
			err = "ERROR-800.35 - After parsing Semester " + " \n " + err    
		    return err
		}		
		
		if (inStr.indx < 0 || inStr.indx >= inStr.len) {
			err = "ERROR-800.36 - Missing Year Data "  + " \n " + err	    
			return err		
		}			
		
		
		err = skipSpacesDelims(inStr)
		if (err != "") {
			err = "ERROR-800.37 - Skipping Spaces searching for Year " + " \n " + err   
		    return err
		}
		
		if (inStr.indx < 0 || inStr.indx >= inStr.len) {
			err = "ERROR-800.38 - Missing Year Data "  + " \n " + err	    
			return err			
		}	
		
		
		err  = getYearToken(inStr, tokenArr)
		if (err != "") {
			err = "ERROR-800.39 - Getting Year Token " + " \n " +  err 
		    return err
		}
				
	} // if (isLetter(char))
	
	if (err != "") {
		err = "ERROR-800.70 - Getting Offer Session " + " \n " + err  
	 }
	

	if (LetsTrace) {
		fmt.Printf("..TRACE-    800.90 : OUT : getClassSession() \n")
	}		
	
	return err
			
}

// Parse and Extract the Year Token for the [OfferSession] Field
// funcid:920
func getYearToken (inStr *ChStr, tokenArr []string ) string {
	 var retToken string
	 // var retYear  string
	 var err      string
	 
	 
	 if (LetsTrace) {
		fmt.Printf("....TRACE-  920.10 : IN- : getYearToken() %v \n", inStr)
	}	
		 
	 
	 retToken, err = getNumberToken(inStr)
	 if (err != "") {
	 	err = "ERROR-920.30 - When Getting Year data " + " \n " + retToken + " \n " + err 
	 	return err
	 }
	 
	 
	 
	 tokenArr[Year], err  = validateYear(retToken)	
	 if (err != "") {
	 	err = "ERROR-920.35 - Invalid Year.  Or Course not offered for Year " + retToken + " \n "	
	 	return err	 	
	 }
	 
	 
	 
	 if (LetsTrace) {
		fmt.Printf("....TRACE-  920.90 : OUT : getYearToken() %v retToken %v tokenArr %v \n", inStr, retToken, tokenArr)
	}	 
	  
	  return ""

}

// Parse and Extract the Semester Token for the [OfferSession] Field
//funcid:950
func getSemesterToken (inStr *ChStr, tokenArr []string ) string {
	 var retToken string
	 var err      string
	 
	 
	 if (LetsTrace) {
		fmt.Printf("....TRACE-  950.10 : IN- : getSemesterToken() %v \n", inStr)
	}	
		  
	 
	 retToken, err = getAlphaToken(inStr)
	 if (err != "") {
	 	err = "ERROR-950.30 - When Getting Semester data " + " \n " + err 
	 	return err
	 }	 
	 
	 tokenArr[Semester], err = validateSemester(strings.ToUpper(retToken))	
	 if (err != "" ) {
	 	err = "ERROR-950.35 - Invalid Semester Entry " + retToken + " \n"
	 	return err	 	
	 }
	 
	 
	 
	  if (LetsTrace) {
		fmt.Printf("....TRACE-  950.90 : OUT : getSemesterToken() %v retToken %v  tokenArr[Semester]-%v\n", inStr, retToken, tokenArr[Semester])	 
	  }
	 
	 return err

}


// Validate Course Offer Year  using a Year Range Validator
//funcid:970
  func validateYear(yearStr string) (string, string) {
  	var err string
  	
  if (LetsTrace) {
		fmt.Printf("....TRACE-  970.10 : IN- : validateYear() \n")
	}		
  	numYear, errGO := strconv.Atoi(yearStr)
  	if (errGO != nil) {
  		err = "ERROR-970.20 - Invalid Year Input" + yearStr + " \n " + err
  		return "", err
  	} // if (errGO != nil) 
  	
  	if (numYear < 99) {
  		numYear += 2000
  	} // fix Year abbreviation
  	
  	if (numYear < EarliestCourseYear) || (numYear >= LatestCourseYear) {
  		err = "ERROR-970.70 - Invalid Year Range " + string(numYear) + " \n " + err
  		return "", err
  	}
  	
  	validYear := strconv.Itoa(numYear)
  	
  if (LetsTrace) {
		fmt.Printf("....TRACE-  970.90 : OUT : validateYear() %v \n", validYear)
	}		
  
  	return validYear, ""
  }


// Validate Semester using a Semester Dictionary 
// funcid:975
func validateSemester(semesterStr string) (string, string) {
	var validSemester string
	var inMap bool
	var err string

	if (LetsTrace) {
		fmt.Printf("....TRACE-  975.10 : IN- : validateSemester() \n")
	} // if (LetsTrace)
    validSemester, inMap = ValidSemester[semesterStr]
    if !(inMap) {
    	err = "ERROR-975.15 - Invalid Semester lookup" + " \n " + err
    	return "", err
    }
    
    if (LetsTrace) {
		fmt.Printf("..TRACE-    975.90 : OUT : validateSemester() %v \n", validSemester)    
    } // if (LetsTrace)
    return validSemester, ""
}


//=====================================================================
// Function main() - Contains Primary Parser for Input String
//=====================================================================

//funcid:1000
func main() {
	 
 // Setup main() scoped variables
 var c   byte
 var err string 
 
 
 // Setup Runtime Trace/Debug Environment
 LetsTrace = true 
 
 //------------------------------------------------------------------
 // -------------CLI Test Harness Code - FORever Loop ---------------
 // -- Comment Out for IDE testing.  Note import section and } at end
 //------------------------------------------------------------------
 reader := bufio.NewReader(os.Stdin)
 fmt.Printf ("\n")
 fmt.Println("Course Selection Entry - Type Quit to Exit")
 fmt.Println("-------------------------------------------")
 fmt.Printf ("\n")
 

 for {
    fmt.Print("-> ")
    inputStr, _ := reader.ReadString('\n')
    // convert CRLF to LF
    inputStr = strings.Replace(inputStr, "\n", "", -1)
    

    if strings.Compare("QUIT", strings.ToUpper(inputStr)) == 0 {
      fmt.Println("Exiting Course Selection")
      break
    } // if strings.Compare("hi", text)

 //------------------------------------------------------------------
 
 
 // Manually set inputStr for IDE testing
 // inputStr := "    CS-111 Fall 2019"
 
 // Setup INPUT Data Structures
 InputStrStruct.data = inputStr
 InputStrStruct.indx = 0
 InputStrStruct.len  = len(inputStr)
 

 // Setup OUTPUT Data Structure 
 tokenList := make([]string, CurTokens, MaxTokens)
 tokenList[Dept]    = ""
 tokenList[Course]  = ""
 tokenList[Semester]= ""
 tokenList[Year]    = "" 
 
 
 
 if (LetsTrace) {
 	fmt.Printf("TRACE-     1000.100: IN  : main().1000.10  InputStrStruct => %v \n", InputStrStruct)
 }
 
 // =====================================================================
 // Skip Leading Spaces and Delimiters
 // ===================================================================== 
	
	if (InputStrStruct.data == "") {
		err = "ERROR-1000.107 No Input Data Found \n"   
	   goto ExitMain		
	}
	
	err = skipSpacesDelims(&InputStrStruct)
 
	if (err != "") {
	   err = "ERROR-1000.150 at start of input string " + InputStrStruct.data + "  \n " +  err 	   
	   goto ExitMain
	} // if (err != "")  
 
 
 // ===================================================================== 
 // Process for Department Course data
 // ===================================================================== 
 err = getDeptCourse (&InputStrStruct, tokenList)
 

 if (err != "") {
	err = "ERROR-1000.200 - in main() " + "  \n " + err
     goto ExitMain
 }	

 
 if (LetsTrace) {
	 fmt.Printf("\n =========================================================================\n")	
 	 fmt.Printf ("\nTRACE-      1000.250 - Result After Parsing [DeptCourse] field %v \n", tokenList)
	 }

// =====================================================================
// Check for Field Seperator between [DeptCourse] & [OfferSession] Fields
// Assumption : Expecting exactly ONE Field Seperator, "a Space", between 
//              [DeptCourse] and {OfferSession] Fields. 
// =====================================================================
 if (LetsTrace) {
	 	fmt.Printf("\n =========================================================================\n\n")
}

 
 
 if (InputStrStruct.indx < 0 || InputStrStruct.indx >= InputStrStruct.len) {
		err = "ERROR-1000.500 - Missing Field Seperator and Session Data " + " \n " + err		
		goto ExitMain			
	}	
  
 c = InputStrStruct.data[InputStrStruct.indx] 
 
 if (LetsTrace) {
	 	fmt.Printf("TRACE-      1000.550: MID : main() Field Seperator [DeptCourse]'%v'[CourseSession] \n", string(c))  
 }
 
 
 if (c != FieldSeperator) {
 	err = "ERROR-1000.550 - Missing Field Seperator between [DeptCourse] and [OfferSession].  Expecting '" + string(FieldSeperator) + "' but finding '" + string(InputStrStruct.data[InputStrStruct.indx]) +" \n " 	
 	goto ExitMain 		
 	
 }
 
 if (InputStrStruct.indx + 1 >= InputStrStruct.len) {
		err = "ERROR-1000.555 - Missing Class Offer Session Data - Year, Semester " + InputStrStruct.data + "\n " + err		
		goto ExitMain		
 }	
 
 // Skipping the Fieled Seperator
 InputStrStruct.indx++
 
 c = InputStrStruct.data[InputStrStruct.indx]
 
 
 if (LetsTrace) {
	 	fmt.Printf("\n =========================================================================\n")
 }
 
 
 // =====================================================================  
 // Continue to Parse next Field for Course Offer Session Data 
 // =====================================================================  
 if (isLetter(c) || isNumber(c)) {
	 err = getOfferSession (&InputStrStruct, tokenList)
     if (err != "") {
     	err = "ERROR-1000.556 - Error while getting [OfferSession] Data \n " + err 
     }     
     goto ExitMain
 }   
 
 // Delimiter "other than and different" from  FieldSeperator found before Offer Session Field
 if (isDelimiter(c)) {
 	err = "ERROR-1000.770 - Only One Delimiter allowed  between [DeptCourse] and [OfferSession]. Found Char ==>'" + string(c) + "'" + "\n " + err 
 	goto ExitMain																			
 }	
 	                                            
  
 // Invalid Character found before Offer Session Field    
 if !(isValid(c)) {
 	 err = "ERROR-1000.850 main() founbd Invalid Char \n" +  err
     	goto ExitMain 	 
	 } // if !(isValid(c))
 
 

 ExitMain:	
	    if (err != "") {
	    	//fmt.Printf("\nInput Entry   |==> [%v]\n", InputStrStruct.data)	    	
	    	// fmt.Printf("InputStrStruct|==> %v\n", InputStrStruct) 	    	
	    	fmt.Printf("\nError STACK   |==> \n-----------------\n[%v]\n-----------------\n", err)
	    }
	    
	     if (LetsTrace) {	    
	     	fmt.Printf("TRACE-     1000.900: OUT : main()") 	     	
	     }

	    // Print FINAL Results
	    if LetsTrace {
	    	fmt.Printf("\n")
	    }
	    fmt.Printf("\nInput Entry   |==> [%v]\n", InputStrStruct.data)	    
	 	fmt.Printf("Output Object |==> %v \n\n",  tokenList)
	 

 } // for  Console Entry Loop

} // main









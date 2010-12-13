package fs

/*
char fsdata[] = { 'o', 'k', 0 };
int fsdataLen = sizeof(fsdata);
*/
import "C"

import "fmt"

func Hi() {
	len := C.fsdataLen
	fmt.Println("len = ", len)

	//println("Slice is: ", C.fsdata)
	//println("chars: ", C.fsdata[0], C.fsdata[1])
	for c := range C.fsdata {
		fmt.Println("Char: ", c)
	}
}

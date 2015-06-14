This project is dead. Here's what it was:

This is a patch that made Go's tiny runtime more capable. With this patch applied on top of your Go source tree, you could catch interrupts, and compile many more of the standard libraries.

In the Go distribution, the Tiny runtime is just a toy to demo how to run Go on raw hardware. Adding more complications to it would obscure the minimum things it is trying to show. So instead, we'll maintain this patch separately.
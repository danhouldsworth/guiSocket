Superceded & ported to C as part of [helloLanguages](https://github.com/danhouldsworth/helloLanguages/tree/master/IO/Networking/GUIsocket) project.

Below for reference and interest only...

# guiSocket
HTML5 Canvas over WebSocket as a GUI frontend for backend computation

#Example
See example/mandyGUI.go

Run with

    cd example
    go run mandyGUI.go

and then you will be prompted to navigate a browser to the IP set (eg ```127.0.0.1:8888``` ).

#API

    guiSocket.Screen(s int)       // Sets the screen width & height. Default 1024.
    guiSocket.Address(s string)   // Set listening address + port. Default 127.0.0.1:8888

    guiSocket.Launch()            // Run with current settings

    guiSocket.Wipe()                    // Code 000 - Clears screen
    guiSocket.Move(x, y)                // Code 001 - todo
    guiSocket.Plot(x, y, C)             // Code 010 - Plot pixel at x,y in colour C
    guiSocket.DrawTo(x, y, C)           // Code 011 - todo
    guiSocket.FillRect(x, y, w, h, C)   // Code 100 - Fill rectanble at x,y of width/height w,h in colour C
    guiSocket.Circle(x, y, r, C)        // Code 101 - Draw circle of raduis r at x,y in colour C
    guiSocket.Image(x, y, w, h, Buff)   // Code 110 - todo
    tbc                                 // Code 111 - reserved

    // x,y,w,h  : are 16-bit unsigned ints / aka cartesian coords from 0-65535
    // C        : is a 32bit colour expressed as 4 arguments of octets in the order RGBA


#Payload

    Byte 1 : 76543210
             xxxxxccc   // where xxxxx is a 5-bit number of commands, and ccc us a 3-bit command code

    Btye 2+ : Arguments as defined above. Note that 16-bit numbers are passed as high-byte, low-byte

Use of ccc for multiple commands of the same type, enables efficient use of bandwidth by saving on Websocket, TCP, IP, Ethernet payloads.


#Customisation

The HTML viewer is parsed and served from ```GUIdisplay.html``` and can be easily customised.


#Future roadmap
Alternatively could use the 5 bits as bit-flags. Say write method : XOR v Overwrite, Relative vs Absolute coords, etc...
We'd only need 1 bit for payload length. If its set, then look to next byte for payload length (2-255) otherwise assume single packet.

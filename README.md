# guiSocket
HTML5 Canvas over WebSocket as a GUI frontend for backend computation

#Example
See example/mandyGUI.go

Run with

    go run mandyGUI.go

and then navigate to ```127.0.0.1:8888``` in browser.

#API
    // X,Y,Width,Height : 16bit coords
    // C : 32bit RGBA colour

    guiSocket.Plot( X, Y, C)
    guiSocket.FillRect( X, Y, Width, Height, C)



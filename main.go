package main

import (
    "fmt"
    "os"
    "strconv"

    "github.com/godbus/dbus/v5"
)

func main() {
    conn, err := dbus.ConnectSessionBus()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
        os.Exit(1)
    }
    
    var brightness int = 0
    var max_brightness int = 0
    defer conn.Close()
    pm := conn.Object("org.kde.Solid.PowerManagement", "/org/kde/Solid/PowerManagement/Actions/BrightnessControl")
    call := pm.Call("org.kde.Solid.PowerManagement.Actions.BrightnessControl.brightness", 0)
    if call.Err != nil || call.Store(&brightness) != nil {
        fmt.Fprintln(os.Stderr, "Failed to get current brightness:", err)
        os.Exit(1)
    }

    call = pm.Call("org.kde.Solid.PowerManagement.Actions.BrightnessControl.brightnessMax", 0)
    if call.Err != nil || call.Store(&max_brightness) != nil {
        fmt.Fprintln(os.Stderr, "Failed to get max brightness:", err)
        os.Exit(1)
    }

    brightness_pct := (100 * brightness) / max_brightness
    original := brightness_pct

    if len(os.Args) < 3 {
        fmt.Printf("Usage : %v <up|down> <value>\n", os.Args[0])
        os.Exit(1)
    }

    val, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Printf("Error, value %v must be an integer", val)
        os.Exit(1)
    }

    switch os.Args[1] {
    case "up":
        brightness_pct += val
    case "down":
        brightness_pct -= val
    case "set":
        brightness_pct = val
    default:
        fmt.Printf("Unknown command %v\n", os.Args[1])
        os.Exit(1)
    }

    if brightness_pct > 100 {
        brightness_pct = 100
    } else if brightness_pct < 0 {
        brightness_pct = 0
    }

    fmt.Printf("Brighness : %v%% -> %v%%\n", original, brightness_pct)
    pm.Call("org.kde.Solid.PowerManagement.Actions.BrightnessControl.setBrightness", 0, (255*brightness_pct) / 100)

}

Requirements:
    1. App should run nginx inside a goroutine.
    2. App should truncate nginx logs inside another goroutine based on a timer or file size.
    3. App should propagate nginx runtime errors back to stdout and/or a log file and exit.

Details:
    Running nginx:
        Goroutine runs nginx in the foreground.
        Runtime errors are checked for and to be acted upon by the goroutine.
        Goroutine will call exit() if any errors are sent to the channel.

    Rotate logs:
        Goroutine will scan nginx configs to search for log paths.
        Rotate logs based on size or amount of time passed.
        Execute nginx log refresh command.

    Signal handling:
        If a kill signal is sent to the app it should handle it appropriately.
        Shutdown nginx and call exit()
    

# `spitool`

Subcommands:


1. `new`: Find novel TLEs

    ```Shell
    spitool new -prev data/2019-12-16.tle -new data/2019-12-17.tle | wc -l
    ```

1. `old`: Find old TLEs

    ```Shell
    spitool old -prev data/2019-12-16.tle -new data/2019-12-17.tle | wc -l
    ```

1. `csv`: TLEs to CSV

    ```Shell
    spitool csv -in data/2019-12-16.tle | head
    ```

1. `vsc`: CSV to TLEs

    ```Shell
    spitool csv -in data/2019-12-16.tle | spitool vsc -in - | head
    ```

1. `prop`: Propagate TLEs 

    ```Shell
    spitool prop -in data/2019-12-16.tle | head
    ```

1. `elements`: Extract some elements from TLEs

    ```Shell
    spitool elements -in data/2019-12-16.tle | head
    ```

1. `sample`: Sample TLEs

    ```Shell
    spitool sample -mod 3 -rem 0 -in data/active.tle | wc -l
    ```
    
1. `tag`: Add a suffix to (trimmed) line 0 if TLEs

    ```Shell
    spitool tag -in data/planet.tle -tag planetlabs | head
    ```

1. `plot`: Generate a basic plot of report distances

   ```Shell
   cat data/active.tle | spibatch | tee reports.json | spitool plot
   ```

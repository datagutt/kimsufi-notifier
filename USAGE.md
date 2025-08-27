# Usage <img src="./assets/bash.svg" width="24">

## List available servers

```
$ kimsufi-notifier list --help
List servers from OVH Eco (including Kimsufi) catalog

Usage:
  kimsufi-notifier list [flags]

Examples:
  kimsufi-notifier list --category kimsufi
  kimsufi-notifier list --category vps
  kimsufi-notifier list --country US --endpoint ovh-us

Flags:
      --category string       category to filter on (allowed values: kimsufi, soyoustart, rise, vps)
  -d, --datacenters strings   datacenter(s) to filter on, comma separated list (known values: bhs, fra, gra, hil, lon, par, rbx, sbg, sgp, syd, vin, waw, ynm, yyz)
  -h, --help                  help for list
  -p, --plan-code string      plan code to filter on (e.g. 24ska01)

Global Flags:
  -c, --country string     country code, known values per endpoints:
                             ovh-eu: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                             ovh-ca: ASIA, AU, CA, IN, QC, SG, WE, WS
                             ovh-us: US
                            (default "FR")
  -e, --endpoint string    OVH API Endpoint (allowed values: ovh-ca, ovh-eu, ovh-us) (default "ovh-eu")
  -l, --log-level string   log level (allowed values: panic, fatal, error, warning, info, debug, trace) (default "error")
```

### VPS Examples

List VPS servers with availability data:
```bash
# List all VPS plans
kimsufi-notifier list --category vps --country FR

# List VPS available in specific datacenters
kimsufi-notifier list --category vps --datacenters GRA,SBG --country FR

# Check specific VPS plan
kimsufi-notifier list --category vps --plan-code vps-starter-1-2-20 --country FR
```

## Check availability

```
$ kimsufi-notifier check --help
Check OVH Eco (including Kimsufi) server availability

datacenters are the available datacenters for this plan

Usage:
  kimsufi-notifier check [flags]

Examples:
  kimsufi-notifier check --plan-code 24ska01
  kimsufi-notifier check --plan-code 24ska01 --datacenters gra,rbx
  kimsufi-notifier check --plan-code vps-starter-1-2-20 --country FR

Flags:
  -d, --datacenters strings     datacenter(s) to filter on, comma separated list (known values: bhs, fra, gra, hil, lon, par, rbx, sbg, sgp, syd, vin, waw, ynm, yyz)
      --help                    help for check
  -h, --human count             Human output, more h makes it better (e.g. -h, -hh)
      --list-options            list available item options
  -o, --option stringToString   options to filter on, comma separated list of key=value, see --list-options for available options (e.g. memory=ram-64g-noecc-2133) (default [])
  -p, --plan-code string        plan code name (e.g. 24ska01, vps-starter-1-2-20)

Global Flags:
  -c, --country string     country code, known values per endpoints:
                             ovh-eu: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                             ovh-ca: ASIA, AU, CA, IN, QC, SG, WE, WS
                             ovh-us: US
                            (default "FR")
  -e, --endpoint string    OVH API Endpoint (allowed values: ovh-ca, ovh-eu, ovh-us) (default "ovh-eu")
  -l, --log-level string   log level (allowed values: panic, fatal, error, warning, info, debug, trace) (default "error")
```

### VPS Check Examples

The check command automatically detects VPS plan codes and provides detailed datacenter availability with OS-specific status:

```bash
# Check VPS availability across all datacenters
kimsufi-notifier check --plan-code vps-starter-1-2-20 --country FR

# Check VPS availability in specific datacenters
kimsufi-notifier check --plan-code vps-2025-model1 --datacenters GRA,SBG --country WE --endpoint ovh-ca

# Example output:
# planCode              datacenter    status       linuxStatus    windowsStatus
# --------              ----------    ------       -----------    -------------
# vps-starter-1-2-20    SBG           available    available      out-of-stock
# vps-starter-1-2-20    GRA           available    available      out-of-stock
```

**VPS vs Dedicated Server Check Differences:**
- **VPS**: Shows datacenter-by-datacenter availability with Linux/Windows status
- **Dedicated Servers**: Shows memory/storage configuration availability

**VPS Plan Code Detection:**
VPS plans are automatically detected by prefixes like:
- `vps-*` (e.g., `vps-starter-1-2-20`, `vps-essential-2-4-40`)
- `s1-*` (e.g., `s1-2` for VPS 2018 SSD 1)

## Order a server

```
$ kimsufi-notifier order --help
Place an order for a servers from OVH Eco (including Kimsufi) catalog

Usage:
  kimsufi-notifier order [flags]

Examples:
  kimsufi-notifier order --plan-code 24ska01 --datacenter rbx --dry-run
  kimsufi-notifier order --plan-code 25skle01 --datacenter bhs --item-option memory=ram-32g-noecc-1333-25skle01,storage=softraid-3x2000sa-25skle01

Flags:
      --auto-pay                            automatically pay the order
  -d, --datacenter string                   datacenter (known values: bhs, fra, gra, hil, lon, par, rbx, sbg, sgp, syd, vin, waw, ynm, yyz)
  -n, --dry-run                             only create a cart and do not submit the order
  -h, --help                                help for order
  -i, --item-configuration stringToString   item configuration, see --list-configurations for available values (e.g. region=europe) (default [])
  -o, --item-option stringToString          item option, see --list-options for available values (e.g. memory=ram-64g-noecc-2133-24ska01) (default [])
      --list-configurations                 list available item configurations
      --list-options                        list available item options
      --list-prices                         list available prices
      --ovh-app-key string                  environement variable name for OVH API application key (default "OVH_APP_KEY")
      --ovh-app-secret string               environement variable name for OVH API application secret (default "OVH_APP_SECRET")
      --ovh-consumer-key string             environement variable name for OVH API consumer key (default "OVH_CONSUMER_KEY")
  -p, --plan-code string                    plan code name (e.g. 24ska01)
      --price-duration string               price duration, see --list-prices for available values (default "P1M")
      --price-mode string                   price mode, see --list-prices for available values (default "default")
  -q, --quantity int                        item quantity (default 1)

Global Flags:
  -c, --country string     country code, known values per endpoints:
                             ovh-eu: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN
                             ovh-ca: ASIA, AU, CA, IN, QC, SG, WE, WS
                             ovh-us: US
                            (default "FR")
  -e, --endpoint string    OVH API Endpoint (allowed values: ovh-ca, ovh-eu, ovh-us) (default "ovh-eu")
  -l, --log-level string   log level (allowed values: panic, fatal, error, warning, info, debug, trace) (default "error")
```


## VPS Support

The tool now supports both OVH Eco dedicated servers (Kimsufi, So you Start, Rise) and VPS instances. VPS support includes:

### Features
- **Real-time availability data**: Live status from OVH VPS availability API
- **Datacenter filtering**: Filter VPS by specific datacenters
- **OS-specific availability**: Separate Linux and Windows availability status
- **Automatic detection**: VPS plans automatically detected by plan code prefix

### VPS Categories
- `vps`: General VPS category covering all VPS plans
- Includes various VPS families like:
  - VPS Starter (e.g., `vps-starter-1-2-20`)
  - VPS Value (e.g., `vps-value-1-2-40`)  
  - VPS Essential (e.g., `vps-essential-2-4-40`)
  - VPS Comfort (e.g., `vps-comfort-4-8-160`)
  - VPS Elite (e.g., `vps-elite-8-16-320`)
  - Legacy VPS 2018 (e.g., `s1-2`)

### Plan Code Detection
VPS plans are automatically detected by these prefixes:
- `vps-*`: Modern VPS plans (Starter, Value, Essential, Comfort, Elite)
- `s1-*`: Legacy VPS 2018 SSD series

### Country and Endpoint Support
VPS availability varies by OVH endpoint:
- **ovh-eu**: Europe (FR, DE, GB, etc.)
- **ovh-ca**: Canada and Asia-Pacific (CA, WE, SG, etc.)  
- **ovh-us**: United States

### Datacenter Codes
Common VPS datacenters include:
- **BHS**: Beauharnois, Canada
- **GRA**: Gravelines, France  
- **SBG**: Strasbourg, France
- **DE**: Germany (Frankfurt)
- **UK**: United Kingdom (London)
- **WAW**: Warsaw, Poland
- **SGP**: Singapore
- **SYD**: Sydney, Australia

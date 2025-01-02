# trend aad

## Debugging with the UI

```bash
go run main.go -- -rod=show
```

## How to use

`trend_aad` will query all aws credentials from your account and store at `~/.aws/credentials`.

### Command parameters

```bash
./trend_aad "trend user" "trend password"
```

### OS environment variable

```bash
export TREND_USERNAME=xxx
export TREND_PASSWORD=xxx
./trend_aad
```

### Terminal interactive

```bash
./trend_aad
```

## Switch AWS profile for fish shell

add a aws-switch profile function to `~/.config/fish/functions/aws-switch.fish`

this script use the tool `fzf`

```fish
function aws-switch
    set profiles (aws configure list-profiles)
    if test (count $profiles) -eq 0
        echo "No AWS profiles found."
        return 1
    end

    set selected (printf '%s\n' $profiles | fzf --height 40% --reverse --border)

    if test -n "$selected"
        set -gx AWS_PROFILE $selected
        echo "Switched to AWS profile: $AWS_PROFILE"
    else
        echo "No profile selected."
    end
end


```

## install chromium dependencies

```bash
make install_chromium_dependencies
```

export GODIR_INSTALL_PATH="$HOME/Dev/godir"
export GODIR_SCRIPT="$GODIR_INSTALL_PATH/main"
export GODIR_DB_FILE_PATH="$HOME/.godir/godir.db"

godir() {
    if [ "$1" == "--flush" ]; then

        $GODIR_SCRIPT --flush

    elif [ "$1" == "--prev" ] || [ "$1" == "-p" ]; then

        prev=$($GODIR_SCRIPT --prev)

        echo "PREV: $prev"

        cd $prev

    elif [ "$1" == "--fuzzy" ] || [ "$1" == "--fzf" ]; then

        if [ ! -e $(which fzf) ]; then
            echo "Sorry, missing FZF"

            exit 1
        fi

        dir=$(cat $GODIR_DB_FILE_PATH | fzf)

        godir "$dir"
    elif [ "$1" == "--show-db" ] || [ "$1" == "-db" ]; then
        echo "DB Contents: "
        
        cat "$GODIR_DB_FILE_PATH"
    else
        dir=$1
        dir="${dir/#~/${HOME}}"

        if [ -e "$1" ]; then
            echo "CDing $1"

            dir="$(readlink -f "$1")"

            echo "ABS PATH: $dir"

            cd $($GODIR_SCRIPT "$dir")
        elif [ "$1" == "-" ]; then
            echo "CDing previous directory"

            cd $($GODIR_SCRIPT "$1")
        else
            echo "Path \"$dir\" not exists"

            # exit 1
        fi
    fi
}

alias gd="godir"
alias gdp="godir --prev"
alias gdf="godir --flush"
alias gdi="godir --fuzzy"
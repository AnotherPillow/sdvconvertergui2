<script>

    let _toastVisible = false;
    let transitionTime = 0;
    let toastContent = ''
    let toastBg = null;
    let toastColour = null;

    export const SendToast = (
        message = 'This is a Toast!',
        time = 2000,
        fade = 100,
        bg = 'wheat',
        colour = 'black'
    ) => {
        _toastVisible = true;
        transitionTime = fade;
        toastContent = message
        toastBg = bg;
        toastColour = colour

        setTimeout(() => {
            _toastVisible = false
        }, time + fade)
    }

</script>

<div id="sdv-style-toast-container">
    <div id="sdv-style-toast" class={`${_toastVisible ? 'visible' : 'hidden'}`} style={
        `--fadeTime: ${transitionTime}ms; ` +
        `--bg: ${toastBg ? toastBg : 'wheat'}; ` +
        `--col: ${toastColour ? toastColour : 'black'}; ` 
    }>
        <h2 class="text-lg font-bold">{toastContent}</h2>
    </div>
</div>

<style lang="scss">
    @use "../data";

    #sdv-style-toast-container {
        min-width: 100vw;
        min-height: 100vh;

        position: absolute;
        inset: 0;
        z-index: 100;

        pointer-events: none;

        #sdv-style-toast {
            $width: 15em;

            display: grid;
            place-items: center;

            position: absolute;
            right: -($width + 2em);
            bottom: 10vh;

            width: $width;
            height: 8em;
            
            background-color: var(--bg);
            @include data.border_image(8px, "src");
            
            // transition: all var(--fadeTime) linear;
            transition: right var(--fadeTime) linear;

            &.visible {
                // display: block;

                right: 1em;
            }

            &.hidden {
                // display: none;

                right: -($width + 2em);

            }

            h2 {
                color: var(--col);
            }
            
        }
    }

</style>
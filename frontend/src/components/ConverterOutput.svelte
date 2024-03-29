<script lang="ts">

    export let converterOutputText = 'Run a converter to see its output...'

    window.runtime.EventsOn('OUTPUT_CONVERTER_TEXT', (text) => {
        converterOutputText += '\n' + text.replace(/\x1B[[0-9;]*[mGKH]/gi, '')
    })

</script>


<textarea id="converter-output" class="code" disabled>{converterOutputText}</textarea>




<style lang="scss">
    @use "../data";

    textarea {
        height: 70%;
        min-height: 2em;
        width: 90%;
        
        margin: 2em 0;
        padding: 5px;
        font-size: 11px;
        color: rgb(62, 62, 62);
        
        scrollbar-width: none;
        cursor: text;
        
        @include data.border_image(8px, "../");
        border-radius: 8px;
        
        background-color: rgb(200, 182, 147);

        outline: none;
        &:active, &:focus, &:focus-within {
            outline: none;

        }
    }
</style>
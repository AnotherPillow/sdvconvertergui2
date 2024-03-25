<script lang="ts">
    import filesBtn from '../assets/images/files_btn.png'
    import { ChooseManifest } from '../../wailsjs/go/main/App.js'
    import { dialogOpen } from '../stores'
    import { get } from 'svelte/store';
    import JSON5 from 'json5'
    import { catchAndNotify } from '../util';

    export let selectedManifest = null;
    export let selectedFilePath = '';

    const onChange = (e) => {
        
    }

    const openDialogForFileSelection = (e) => {
        console.log('clicked!')
        e.preventDefault();
        ChooseManifest().then(data => {
            if (data.content) data.content = catchAndNotify(
                ()=>JSON5.parse(data.content.replace(/[\r\n]/g, '')), 'Unreadable manifest selected!')
            console.log(data)
            if (!data.filename) return;
            selectedFilePath = data.filename

            if (!data.content.Name || !data.content.Author || !data.content.UniqueID || !data.content.Description) {
                window.SendToast("Invalid manifest selected!", 2000, 100, '#f5bbb3')
                return selectedManifest = null;
            } else if (!data.content?.ContentPackFor?.UniqueID) {
                window.SendToast("C# mod selected - please choose a content pack.", 5000, 100, '#f5bbb3')
                return selectedManifest = null;
            }

            selectedManifest = data
        })
    }
    let _file_input: HTMLElement | null = null;
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<span id="file-uploader-area" class="flex" on:click={openDialogForFileSelection}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <label for="file-input" class="file-label" >
        <img 
            src={filesBtn}
            alt="Upload"
            class="file-btn clickable"
            width="50"
            height="50"
        />
    </label>
    <input 
        type="file"
        class="file-input"
        id="file-input" 
        accept=".json"
        on:change={onChange}
        bind:this={_file_input}
    />
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <h3 on:click={() => {
        // _file_input?.click()
    }}>
        {#if selectedManifest == null}
            Choose a manifest.json
        {:else}
            {selectedManifest.content.Name}
        {/if}
    </h3>
</span>

<style lang="scss">

    #file-uploader-area {
        cursor: pointer;
        width: max-content;
        color: var(--main-color);

        align-items: center;
        font-weight: bold;

        #file-input {
            opacity: 0;
            width: 0;
        }

        h3 {
            padding-left: 5px;
        }
    }

</style>
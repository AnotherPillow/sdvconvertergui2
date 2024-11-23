<script lang="ts">
    // @ts-ignore
    import Check from "lucide-svelte/icons/check";
    // @ts-ignore
    import ChevronsUpDown from "lucide-svelte/icons/chevrons-up-down";
    import * as Command from "../lib/components/ui/command";
    import * as Popover from "../lib/components/ui/popover";
    import { Button } from "../lib/components/ui/button";
    import { cn } from "$lib/utils";
    import { tick } from "svelte";
    import { converters } from '../stores'
    import dropdownArrow from '../assets/images/dropdown-arrow.png'
    import repoArrow from '../assets/images/repo_arrow.png'
    import { get } from "svelte/store";
   
    let frameworks = [];
    converters.subscribe(c => {
        if (c.length == 0) return;
        frameworks = c
    })

    export let selectedValue = ''
   
    let open = false;
    let value = "";

    $: {
        console.log(frameworks)
        selectedValue =
            frameworks.find((f) => f.Name === value)?.Name ??
            "Choose a Converter...";
    }
   
    // We want to refocus the trigger button when the user selects
    // an item from the list so users can continue navigating the
    // rest of the form with the keyboard.
    function closeAndFocusTrigger(triggerId: string) {
        open = false;
        tick().then(() => {
            document.getElementById(triggerId)?.focus();
        });
    }
</script>
   
<Popover.Root bind:open let:ids>
    <Popover.Trigger asChild let:builder>
        <Button
            builders={[builder]}
            variant="primary"
            role="combobox"
            aria-expanded={open}
            class="justify-between py-0 sdv-dropdown relative rounded-none sdv px-2 pr-6"
        >
            {selectedValue}
            <!-- <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" /> -->
            <img src={dropdownArrow} alt="arrow" class="z-10 absolute right-0 sdv-dropdown-arrow"/>
            {#if selectedValue != '' && 
                frameworks.find(x => x.Name == selectedValue)?.Repo != undefined}
                <span id="show-converter-repo">
                    <a href={frameworks.find(x => x.Name == selectedValue)?.Repo} target="_blank" title="Open Repository">
                        <img src={repoArrow} alt=">">
                    </a>
                </span>
            {/if}
        </Button>
    </Popover.Trigger>
    <Popover.Content class="sdv-dropdown-content-outer p-0">
        <Command.Root class="sdv-dropdown-content">
            <Command.Input searchVisible={false} 
                class="sdv text-md text-black converter-searchbox" 
                oClass="searchbox-sdv-seperator-bottom" 
                placeholder="Search Converters..." />
            <Command.Empty>No framework found.</Command.Empty>
            <Command.Group class="sdv">
                {#each frameworks as framework}
                    <Command.Item
                        value={framework.Name}
                        class={`text-lg pl-4 ${framework.Name == selectedValue ? 'bg-sdv_dropdown_unsel' : ''}`}
                        onSelect={(currentValue) => {
                            value = currentValue;
                            closeAndFocusTrigger(ids.trigger);
                        }}
                    >
                        <!-- <Check
                            class={cn(
                            "mr-2 h-4 w-4",
                            value !== framework.value && "text-transparent"
                            )}
                        /> -->
                        {framework.Name}
                    </Command.Item>
                {/each}
            </Command.Group>
        </Command.Root>
    </Popover.Content>
</Popover.Root>



<style lang="scss">
    @use "../data";

    img {
        background-image: data.$sdv-dropdown-img;
    }

    $dropdown-width: 200px;

    $dropdown-height: calc(#{$dropdown-width} / 6);

    :global(.sdv-dropdown) {
        @include data.dropdownBorder;
        @include data.dropDownBorderRounding;
        
        min-width: $dropdown-width;
        width: fit-content;
        height: $dropdown-height;

        font-size: 1.25em;
        
        color: data.$sdv-dropdown-colour;
        background-color: data.$sdv-dropdown-unselected;

        &:hover, &:active, &:focus, &:focus-visible, &:focus-within, &:target {
            @include data.dropdownBorder;
            @include data.dropDownBorderRounding;
            
            box-shadow: none !important;

            filter: brightness(0.8);
        }
    }

    :global(.sdv-dropdown-content) {
        @include data.dropdownBorder;
        border-radius: 0;

        box-shadow: none !important;
        outline:  none !important;
        
        color: data.$sdv-dropdown-colour;
        background-color: data.$sdv-dropdown-unselected;
    }

    :global(.sdv-dropdown-content-outer) {
        border-radius: 0;
        width: calc(#{$dropdown-width} / 1.2);
        
        color: data.$sdv-dropdown-colour;
        background-color: data.$sdv-dropdown-unselected;

        box-shadow: none !important;
        outline:  none !important;
    }

    :global(.converter-searchbox::placeholder) {
        color: black;
        opacity: 0.8;
    }

    .sdv-dropdown-arrow {
        height: inherit;
        aspect-ratio: 1;

        right: calc(#{$dropdown-width} / -15)
    }

    #show-converter-repo {
        position: absolute;
        right: -45px;

        
        
        height: inherit;
        a {
            height: inherit;
            display: flex;
            place-items: center;

            img {
                height: 80%;
            }
        }

    }
</style>
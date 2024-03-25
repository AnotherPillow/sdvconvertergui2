<script>
	import { Command as CommandPrimitive } from "cmdk-sv";
	import { cn } from "$lib/utils";
	export let asChild = false;
	let className = undefined;
	export { className as class };
</script>

<CommandPrimitive.Item
	{asChild}
    aria-selected="false"
	class={cn(
		"aria-selected:bg-sdv_dropdown_unsel relative flex cursor-default select-none items-center  px-2 py-1.5 text-sm outline-none  aria-selected:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",// rounded-sm
		className
	)}
	{...$$restProps}
	let:action
	let:attrs
>
	<slot {action} {attrs} />
</CommandPrimitive.Item>

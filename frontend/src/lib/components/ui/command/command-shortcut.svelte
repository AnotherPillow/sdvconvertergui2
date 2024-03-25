<script>
	import { cn } from "$lib/utils";
	let className = undefined;
	export { className as class };
</script>

<span
	class={cn("ml-auto text-xs tracking-widest text-muted-foreground", className)}
	{...$$restProps}
>
	<slot />
</span>

<script>
    export let center = false;

    let classes = "";
    export { classes as class }; // export reserved keyword
</script>

<div class="page-wrapper {classes}" class:center-content={center}>
    <main class="page-content">
        <slot />
    </main>

    <footer class="page-footer">
        
    </footer>
</div>

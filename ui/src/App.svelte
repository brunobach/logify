<script>
    import "./scss/main.scss";
    import Router  from "svelte-spa-router";
    import routes from "./routes";
    import CommonHelper from "@/utils/CommonHelper";
    import { pageTitle, appName, hideControls } from "@/stores/app";
</script>

<svelte:head>
    <title>{CommonHelper.joinNonEmpty([$pageTitle, $appName, "Logify"], " - ")}</title>
</svelte:head>

<div class="app-layout">
    <div class="app-body">
        <Router {routes} />
    </div>
</div>


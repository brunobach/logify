import { replace }           from "svelte-spa-router";
import { wrap }              from "svelte-spa-router/wrap";
import PageLogs              from "@/components/logs/PageLogs.svelte";


const baseConditions = [
    async (details) => {
        const realQueryParams = new URLSearchParams(window.location.search);

        if (details.location !== "/" && realQueryParams.has(import.meta.env.PB_INSTALLER_PARAM)) {
            return replace("/")
        }

        return true
    }
];

const routes = {
    "*": wrap({
        component: PageLogs,
        conditions: baseConditions.concat([]),
        userData: { showAppSidebar: false },
    }),
};

export default routes;

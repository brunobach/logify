import PocketBase from "pocketbase";

const client = import.meta.env.PB_ENV === "dev" ? new PocketBase(
    import.meta.env.PB_BACKEND_URL,
) : new PocketBase(window.location.href);


export default client;

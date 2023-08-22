import PocketBase from "pocketbase";

const client = new PocketBase(
    import.meta.env.PB_BACKEND_URL,
);


export default client;

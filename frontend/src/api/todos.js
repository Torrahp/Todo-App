import client from "./client";

export async function getTodos() {
    const res = await client.get("/todos");
    return res.data;
}

export async function createTodo(title){
    const res = await client.post("/todos", {title});
    return res.data;
}

export async function updateTodo(id, changes){
    const res = await client.patch(`/todos/${id}`, changes);
    return res.data;
}

export async function deleteTodo(id) {
    await client.delete(`/todos/${id}`);
}
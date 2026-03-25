export interface Task {
    id: string;
    name: string;
    completed: boolean;
}

export interface ApiTask {
    id: number;
    title: string;
    status: number;
    deleted: boolean;
}

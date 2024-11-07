import Login from '../page/login';
import Chat from '../page/chat';

const routers = [
    { path : '/', component : <Login /> },
    { path : '/chat', component : <Chat />}
];

export default routers;


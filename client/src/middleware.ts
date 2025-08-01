import { defineMiddleware } from "astro:middleware";

const UNPROTECTED_ROUTES = [
    /\/login($|\/.*)/,
    /\/create-account($|\/.*)/,
    /\/test($|\/.*)/,
    /\/500($|\/.*)/,
    /\/400($|\/.*)/,
];

const isSafeRoute = (path: string): boolean => {
    return UNPROTECTED_ROUTES.some((pattern) => pattern.test(path));
};

export const onRequest = defineMiddleware(async (context, next) => {
    if (isSafeRoute(context.url.toString())) {
        return next();
    }

    const sessionToken = context.cookies.get('session-token')?.value;


    const myHeaders = new Headers();
    if (sessionToken) {
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Cookie", 'session-token=' + sessionToken)
    }

    const response = await fetch(new URL('http://localhost:8081/get-user'), {
        method: 'GET',
        credentials: 'include',
        headers: myHeaders,
    });

    // let content = undefined;
    // const content = await response.json(); 

    if (!response.ok) {
        context.request.headers.set('x-redirect-to', context.url.pathname);
        return context.redirect(context.url.origin +  '/login');
    }

    return next();
});

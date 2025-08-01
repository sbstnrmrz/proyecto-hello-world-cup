import { defineMiddleware } from "astro:middleware";

export const onRequest = defineMiddleware((context, next) => {
    console.log('\n\n middleware');
    const sessionToken = context.cookies.get('session_token')?.value;
    const UNPROTECTED_ROUTES = [
        /^\/login($|\/.*)/, // matches /login and any query params that may be included
        /^\/test($|\/.*)/, // matches /login and any query params that may be included
        /^\/signup($|\/.*)/,
        /^\/invite($|\/.*)/,
        /^\/hello/,
        /^\/500($|\/.*)/,
        /^\/400($|\/.*)/,
    ];    
    const isSafeRoute = (path: string): boolean => {
        return UNPROTECTED_ROUTES.some((pattern) => pattern.test(path));
    };

    if (isSafeRoute(context.url.toString())) {
        return next();
    }

    if (!sessionToken) {
        return context.redirect(context.url.origin +  "/test");
        // context.url = new URL(context.url.origin +  "/test");
        // return next(new Request(new URL(context.url.origin +  "/test"), {
        //     headers: {
        //         "x-redirect-to": context.url.pathname
        //     }
        // }));
    }

    return next();
});

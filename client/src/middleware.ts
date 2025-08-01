import { defineMiddleware } from "astro:middleware";

const UNPROTECTED_ROUTES = [
    /\/login($|\/.*)/,
    /\/create-account($|\/.*)/,
    /\/home($|\/.*)/,
    /\/profile-student($|\/.*)/,
    /\/profile-teacher($|\/.*)/,
    /\/signup($|\/.*)/,
    /\/test($|\/.*)/,
    /\/500($|\/.*)/,
    /\/400($|\/.*)/,
];

const isSafeRoute = (path: string): boolean => {
    return UNPROTECTED_ROUTES.some((pattern) => pattern.test(path));
};

export const onRequest = defineMiddleware((context, next) => {
    if (isSafeRoute(context.url.toString())) {
        return next();
    }

    const sessionToken = context.cookies.get('session-token')?.value;

    console.log(`sessionToken = ${sessionToken}`)

    if (!sessionToken) {
        context.request.headers.set('x-redirect-to', context.url.pathname);
        return context.redirect(context.url.origin +  '/login');
    }

    return next();
});

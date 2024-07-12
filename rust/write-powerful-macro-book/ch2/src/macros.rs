macro_rules! compose {
    ($last:expr) => {
        $last
    };
    ($head:expr, $($tail:expr),+) => {
        $crate::macros::compose_two($head, compose!($($tail),+))
    };
}

pub(crate) use compose;

pub(crate) fn compose_two<T1, T2, T3, F, G>(f: F, g: G) -> impl Fn(T1) -> T3
where
    F: Fn(T1) -> T2,
    G: Fn(T2) -> T3,
{
    move |x| g(f(x))
}

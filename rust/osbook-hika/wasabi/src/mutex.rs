use crate::result::Result;
use core::{
    cell::SyncUnsafeCell,
    fmt,
    ops::{Deref, DerefMut},
    panic::Location,
    sync::atomic::{AtomicBool, AtomicU32, Ordering},
};

pub struct MutexGuard<'a, T> {
    mutex: &'a Mutex<T>,
    data: &'a mut T,
    location: Location<'a>,
}
impl<'a, T> MutexGuard<'a, T> {
    #[track_caller]
    unsafe fn new(mutex: &'a Mutex<T>, data: &SyncUnsafeCell<T>) -> Self {
        Self {
            mutex,
            data: &mut *data.get(),
            location: *Location::caller(),
        }
    }
}

unsafe impl<'a, T> Sync for MutexGuard<'a, T> {}

impl<'a, T> Deref for MutexGuard<'a, T> {
    type Target = T;

    fn deref(&self) -> &Self::Target {
        self.data
    }
}
impl<'a, T> DerefMut for MutexGuard<'a, T> {
    fn deref_mut(&mut self) -> &mut Self::Target {
        self.data
    }
}

impl<'a, T> Drop for MutexGuard<'a, T> {
    fn drop(&mut self) {
        self.mutex.is_taken.store(false, Ordering::Release)
    }
}

impl<'a, T> fmt::Debug for MutexGuard<'a, T> {
    fn fmt(&self, f: &mut core::fmt::Formatter<'_>) -> core::fmt::Result {
        write!(f, "MutexGuard {{ location: {:?} }}", self.location)
    }
}

pub struct Mutex<T> {
    data: SyncUnsafeCell<T>,
    is_taken: AtomicBool,
    taker_line_num: AtomicU32,
    created_at_file: &'static str,
    created_at_line: u32,
}

impl<T: Sized> fmt::Debug for Mutex<T> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "Mutex @ {}:{}",
            self.created_at_file, self.created_at_line
        )
    }
}

impl<T: Sized> Mutex<T> {
    #[track_caller]
    pub const fn new(data: T) -> Self {
        Self {
            data: SyncUnsafeCell::new(data),
            is_taken: AtomicBool::new(false),
            taker_line_num: AtomicU32::new(0),
            created_at_file: Location::caller().file(),
            created_at_line: Location::caller().line(),
        }
    }
    #[track_caller]
    pub fn try_lock(&self) -> Result<MutexGuard<T>> {
        if self
            .is_taken
            .compare_exchange(false, true, Ordering::Acquire, Ordering::Relaxed)
            .is_ok()
        {
            self.taker_line_num
                .store(Location::caller().line(), Ordering::Relaxed);
            Ok(unsafe { MutexGuard::new(self, &self.data) })
        } else {
            Err("Lock failed")
        }
    }
    #[track_caller]
    pub fn lock(&self) -> MutexGuard<T> {
        for _ in 0..10000 {
            if let Ok(locked) = self.try_lock() {
                return locked;
            }
        }
        panic!(
            "Failed to lock Mutex at {}:{}, caller: {:?}, taker_line_num: {}",
            self.created_at_file,
            self.created_at_line,
            Location::caller(),
            self.taker_line_num.load(Ordering::Relaxed),
        )
    }
    pub fn under_locked<R: Sized>(&self, f: &dyn Fn(&mut T) -> Result<R>) -> Result<R> {
        let mut locked = self.lock();
        f(&mut *locked)
    }
}

unsafe impl<T> Sync for Mutex<T> {}
impl<T: Default> Default for Mutex<T> {
    #[track_caller]
    fn default() -> Self {
        Self::new(T::default())
    }
}

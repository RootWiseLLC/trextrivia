import { Injectable } from '@angular/core';

@Injectable({
    providedIn: 'root'
})
export class SoundService {
    private soundEnabled: boolean = true;
    private readonly STORAGE_KEY = 'jeopardy_sound_enabled';

    constructor() {
        // Load sound preference from localStorage
        const stored = localStorage.getItem(this.STORAGE_KEY);
        if (stored !== null) {
            this.soundEnabled = stored === 'true';
        }
    }

    toggleSound(): void {
        this.soundEnabled = !this.soundEnabled;
        localStorage.setItem(this.STORAGE_KEY, String(this.soundEnabled));
    }

    isSoundEnabled(): boolean {
        return this.soundEnabled;
    }

    setSoundEnabled(enabled: boolean): void {
        this.soundEnabled = enabled;
        localStorage.setItem(this.STORAGE_KEY, String(enabled));
    }

    /**
     * Play audio if sound is enabled
     * @param audio HTMLAudioElement to play
     * @returns Promise that resolves when playback starts or rejects if sound is disabled
     */
    playAudio(audio: HTMLAudioElement): Promise<void> {
        if (this.soundEnabled) {
            return audio.play().catch(err => {
                console.log('Audio play failed:', err);
                return Promise.reject(err);
            });
        }
        return Promise.resolve();
    }
}

import { Component, OnInit, OnDestroy, Input } from '@angular/core';
import { GameStateService } from 'src/app/services/game-state.service';
import { ModalService } from 'src/app/services/modal.service';
import { SoundService } from 'src/app/services/sound.service';

@Component({
    selector: 'app-answer-feedback',
    templateUrl: './answer-feedback.component.html',
    styleUrls: ['./answer-feedback.component.less']
})
export class AnswerFeedbackComponent implements OnInit, OnDestroy {
    private correctAudio: HTMLAudioElement;
    private wrongAudio: HTMLAudioElement;
    private dailyDoubleAudio: HTMLAudioElement;
    private hideTimeout: any;

    constructor(
        protected game: GameStateService,
        protected modal: ModalService,
        private sound: SoundService
    ) {
        this.correctAudio = new Audio('../../../assets/correct.mp3');
        this.wrongAudio = new Audio('../../../assets/wrong.mp3');
        this.dailyDoubleAudio = new Audio('../../../assets/daily_double.mp3');
    }

    ngOnInit() {
        // Play sound effect based on feedback type
        this.playSound();

        // Auto-hide after 1 second (or 3 seconds for birthday)
        const hideDelay = this.isBirthday() ? 3000 : 1000;
        this.hideTimeout = setTimeout(() => {
            this.modal.hideAnswerFeedback();
        }, hideDelay);
    }

    ngOnDestroy() {
        if (this.hideTimeout) {
            clearTimeout(this.hideTimeout);
        }
    }

    playSound() {
        // Check if sound is enabled before playing
        if (!this.sound.isSoundEnabled()) {
            return;
        }

        if (this.isDailyDouble()) {
            this.dailyDoubleAudio.play().catch(err => console.log('Audio play failed:', err));
        } else if (this.isCorrect()) {
            this.correctAudio.play().catch(err => console.log('Audio play failed:', err));
        } else if (!this.isBirthday()) {
            this.wrongAudio.play().catch(err => console.log('Audio play failed:', err));
        }
        // Birthday message plays no sound
    }

    isCorrect(): boolean {
        const feedbackType = this.modal.getAnswerFeedbackType();
        return feedbackType === 'correct';
    }

    isIncorrect(): boolean {
        const feedbackType = this.modal.getAnswerFeedbackType();
        return feedbackType === 'incorrect';
    }

    isTimeout(): boolean {
        const feedbackType = this.modal.getAnswerFeedbackType();
        return feedbackType === 'timeout';
    }

    isDailyDouble(): boolean {
        const feedbackType = this.modal.getAnswerFeedbackType();
        return feedbackType === 'daily-double';
    }

    isBirthday(): boolean {
        const feedbackType = this.modal.getAnswerFeedbackType();
        return feedbackType === 'birthday';
    }
}

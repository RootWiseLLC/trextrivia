import { Component, OnDestroy, OnInit } from '@angular/core';
import { GameStateService } from 'src/app/services/game-state.service';
import { SoundService } from 'src/app/services/sound.service';
import { Question } from 'src/app/model/model';

@Component({
	selector: 'app-board-intro',
	templateUrl: './board-intro.component.html',
	styleUrls: ['./board-intro.component.less']
})
export class BoardIntroComponent implements OnInit, OnDestroy {
	categories: string[];
	questionRows: Question[][];
	private boardIntroAudio: HTMLAudioElement | null = null;
	private soundToggleListener: ((event: Event) => void) | null = null;

	constructor(
		protected game: GameStateService,
		private sound: SoundService,
	) {
		this.categories = this.game.Categories()
		this.questionRows = this.game.QuestionRows()

		try {
			// Randomly select one of the 5 songs
			const songNumber = Math.floor(Math.random() * 5) + 1;
			console.log('[BoardIntro] Selected song:', songNumber);
			this.boardIntroAudio = new Audio(`../../assets/music/song_${songNumber}.mp3`);
		} catch (err) {
			console.error('Failed to initialize board intro audio:', err);
		}
	}

	ngOnInit() {
		console.log('[BoardIntro] Component initialized, sound enabled:', this.sound.isSoundEnabled());

		// Play board intro music if sound is enabled
		if (this.boardIntroAudio && this.sound.isSoundEnabled()) {
			console.log('[BoardIntro] Playing music...');
			this.boardIntroAudio.play().catch(err => console.log('Board intro audio play failed:', err));
		}

		// Listen for sound toggle events
		this.soundToggleListener = (event: Event) => {
			const customEvent = event as CustomEvent;
			const enabled = customEvent.detail.enabled;
			console.log('[BoardIntro] Sound toggled:', enabled);

			if (this.boardIntroAudio) {
				if (enabled) {
					this.boardIntroAudio.play().catch(err => console.log('Board intro audio play failed:', err));
				} else {
					this.boardIntroAudio.pause();
				}
			}
		};
		window.addEventListener('soundToggle', this.soundToggleListener);

		let questionValues = document.getElementsByClassName('question-value') as HTMLCollectionOf<HTMLElement>
		let indexes = Array.from({ length: 30 }, (v, i) => i)
		let revealGroups: any[] = []
		for (let i = 0; i < 6; i++) {
			let revealGroup = []
			for (let j = 0; j < 5; j++) {
				let index = Math.floor(Math.random() * indexes.length);
				let num = indexes.splice(index, 1)[0];
				revealGroup.push(num)
			}
			revealGroups.push(revealGroup)
		}
		setTimeout(() => {
			let i = 0
			let valuesInterval = setInterval(() => {
				if (i < revealGroups.length) {
					for (let j = 0; j < revealGroups[i].length; j++) {
						questionValues[revealGroups[i][j]].style.color = 'var(--jeopardy-yellow)'
					}
					i++
				} else {
					clearInterval(valuesInterval)
				}
			}, 1000)
		}, 2000)
		setTimeout(() => {
			let categoryTitles = document.getElementsByClassName('category-title') as HTMLCollectionOf<HTMLElement>
			let j = 0
			let titlesInterval = setInterval(() => {
				if (j < categoryTitles.length) {
					(categoryTitles[j].firstChild as HTMLElement).classList.add('fade-in')
					j++
				} else {
					clearInterval(titlesInterval)
				}
			}, 2500)
		}, 7000)
	}

	ngOnDestroy() {
		// Stop music when leaving board intro
		if (this.boardIntroAudio) {
			this.boardIntroAudio.pause();
		}

		// Remove event listener
		if (this.soundToggleListener) {
			window.removeEventListener('soundToggle', this.soundToggleListener);
		}
	}
}

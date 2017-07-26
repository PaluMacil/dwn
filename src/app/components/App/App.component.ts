import { Component } from '@angular/core';
import { AccountAPI } from '../../services/AccountAPI.service'

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'app';

  constructor(public accountAPI: AccountAPI){}

  logout() {
    this.accountAPI.Logout().subscribe(
      () => { console.log('Logged out') },
      (err) => { console.log('Could not log out. ' + err.message) }
    );
  }
}

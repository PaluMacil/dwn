// angular
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
// modules
import { AlertModule } from 'ngx-bootstrap';
// services
import { AccountAPI } from '../services/AccountAPI.service'
// pages
// components
import { AppComponent } from '../components/App/App.component';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    AlertModule.forRoot(),
    BrowserModule
  ],
  providers: [
    AccountAPI
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }

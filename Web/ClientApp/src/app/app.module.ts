import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap'; // https://ng-bootstrap.github.io
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome'; // https://github.com/FortAwesome/angular-fontawesome

import { RoutingModule } from './routing/routing.module';

import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { PageComponent } from './page/page.component';


@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    PageComponent
  ],
  imports: [
    BrowserModule,
    NgbModule.forRoot(),
    FontAwesomeModule,

    RoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }

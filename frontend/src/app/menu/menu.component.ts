import { Component, OnInit } from '@angular/core';
import { KeycloakService } from 'keycloak-angular';
import { MenuItem } from 'primeng/api';
import { AppRoutingModule } from '../app-routing.module';
import { Router } from '@angular/router';


@Component({
    selector: 'app-menu',
    templateUrl: './menu.component.html',
    styleUrl: './menu.component.css'
})
export class MenuComponent implements OnInit {
    items: MenuItem[] = [];
    isLoggedIn: boolean = false;
    constructor(private kc: KeycloakService, 
        private router: Router) {

    }
    public async ngOnInit() {
        var isLoggedIn = await this.kc.isLoggedIn();
        var _localItems: MenuItem[] = [];
        if (isLoggedIn) {
            var roles = await this.kc.getUserRoles(true);
        }else{
            return;
        }
        
        this.router.config.forEach(route => {
            if(route.data && route.data['roles'] && route.data['roles'].length > 0){
                //prüfen, ob die Route erlaubt ist anhand der Rolle
                if( ! route.data['roles'].some((role: string)=> roles.includes(role))){
                    return;
                }
            }
            if( route.data && route.data['label']){
                _localItems.push({
                    label: route.data['label'],
                    command: (click)=>{this.router.navigate([route.path])},
                })
            }
        });
        //Es muss das Array überschrieben werden, damit Angular merkt, dass sich etwas geändert hat.
        //Das liegt daran, dass push ein Element hinzufügt, aber die Referenz auf das Array bleibt identisch
        //Erst bei einer neu Zuweisung ändert sich die Referenz und PrimeNG rendert die Menubar neu
        this.items = _localItems;
    }
    login(){
        this.kc.login();
    }
    logout(){
        this.kc.logout();
    }
}

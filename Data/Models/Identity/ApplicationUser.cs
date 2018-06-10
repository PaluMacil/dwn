using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Dwn.Data.Models.Authoring;
using Microsoft.AspNetCore.Identity;

namespace Dwn.Data.Models.Identity
{
    // Add profile data for application users by adding properties to the ApplicationUser class
    public class ApplicationUser : IdentityUser<Guid>
    {
        public string FirstName { get; set; }
        public string LastName { get; set; }
        public long? FacebookId { get; set; }
        public string PictureUrl { get; set; }
        public virtual List<AuthoredContent> AuthoredContent { get; set; }
    }
}
